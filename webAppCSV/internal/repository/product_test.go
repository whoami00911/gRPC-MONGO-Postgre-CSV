package repository

import (
	"database/sql"
	"errors"
	"log"
	"testing"
	"webApp/domain"
	logger "webApp/pkg/logger"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestAddProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	logger := logger.GetLogger()

	product := &ProductRepo{
		db:     db,
		logger: logger,
		dbname: "assets",
	}

	repo := Repository{
		Product: product,
	}

	type mockBehavior func(product domain.Product)

	testTables := []struct {
		name         string
		mockBehavior mockBehavior
		product      domain.Product
		isErr        bool
	}{
		{
			name: "valid",
			mockBehavior: func(product domain.Product) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO \"assets\" \(name, price\) VALUES \(\$1, \$2\)`).WithArgs(product.Name, product.Price).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			product: domain.Product{
				Name:  "name",
				Price: decimal.RequireFromString("50.00"),
			},
			isErr: false,
		},
		{
			name: "Error with transaction",
			mockBehavior: func(product domain.Product) {
				mock.ExpectBegin().WillReturnError(err)
				mock.ExpectRollback()
			},
			isErr: true,
		},
		{
			name: "Exec SELECT Error",
			mockBehavior: func(product domain.Product) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO \"assets\" \(name, price\) VALUES \(\$1, \$2\)`).WithArgs(product.Name, product.Price).WillReturnError(err)
				mock.ExpectRollback()
			},
			product: domain.Product{
				Name:  "name",
				Price: decimal.RequireFromString("50.00"),
			},
			isErr: true,
		},
		{
			name: "Product Exists Error",
			mockBehavior: func(product domain.Product) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO \"assets\" \(name, price\) VALUES \(\$1, \$2\)`).WithArgs(product.Name, product.Price).WillReturnError(domain.ErrProductExists)
				mock.ExpectRollback()
			},
			product: domain.Product{
				Name:  "name",
				Price: decimal.RequireFromString("50.00"),
			},
			isErr: true,
		},
	}

	for _, table := range testTables {
		t.Run(table.name, func(t *testing.T) {
			table.mockBehavior(table.product)

			err := repo.AddProduct(table.product)
			if table.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	logger := logger.GetLogger()

	product := &ProductRepo{
		db:     db,
		logger: logger,
		dbname: "assets",
	}

	repo := Repository{
		Product: product,
	}

	type mockBehavior func(id int, product domain.Product)

	testTables := []struct {
		name         string
		mockBehavior mockBehavior
		product      domain.Product
		id           int
		isErr        bool
	}{
		{
			name: "valid",
			mockBehavior: func(id int, product domain.Product) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id", "name", "price"}).AddRow(product.Id, product.Name, product.Price.StringFixed(2))
				mock.ExpectQuery(`SELECT \"id\", \"name\", \"price\" FROM \"assets\" WHERE \"id\" = \$1`).WithArgs(id).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			id: 1,
			product: domain.Product{
				Id:    1,
				Name:  "name",
				Price: decimal.RequireFromString("50.00"),
			},
			isErr: false,
		},
		{
			name: "Error with transaction",
			mockBehavior: func(id int, product domain.Product) {
				mock.ExpectBegin().WillReturnError(err)
				mock.ExpectRollback()
			},
			isErr: true,
		},
		{
			name: "Product not found",
			mockBehavior: func(id int, product domain.Product) {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT \"id\", \"name\", \"price\" WHERE \"id\" = \$1`).WithArgs(id).WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			isErr: true,
		},

		{
			name: "Scan error",
			mockBehavior: func(id int, product domain.Product) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id", "name", "price"}).
					AddRow(1, "foo", "string")
				mock.
					ExpectQuery(`SELECT \"id\", \"name\", \"price\" WHERE \"id\" = \$1`).
					WithArgs(1).
					WillReturnRows(rows)
				mock.ExpectRollback()
			},
			isErr: true,
		},
	}

	for _, table := range testTables {
		t.Run(table.name, func(t *testing.T) {
			table.mockBehavior(table.id, table.product)

			got, err := repo.GetProduct(table.id)
			if table.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, table.product, got)
			}
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	logger := logger.GetLogger()

	product := &ProductRepo{
		db:     db,
		logger: logger,
		dbname: "assets",
	}

	repo := Repository{
		Product: product,
	}

	type mockBehavior func(id int)

	testTables := []struct {
		name         string
		mockBehavior mockBehavior
		id           int
		isErr        bool
	}{
		{
			name: "Valid",
			mockBehavior: func(id int) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM \"assets\" WHERE \"id\"=\$1`).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			id:    1,
			isErr: false,
		},
		{
			name: "Transaction error",
			mockBehavior: func(id int) {
				mock.ExpectBegin().WillReturnError(err)
				mock.ExpectRollback()
			},
			isErr: true,
		},
		{
			name: "Exec DELETE error",
			mockBehavior: func(id int) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM \"assets\" WHERE \"id\"=\$1`).WithArgs(id).WillReturnError(err)
				mock.ExpectRollback()
			},
			id:    1,
			isErr: true,
		},
		{
			name: "No Rows Affected",
			mockBehavior: func(id int) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM \"assets\" WHERE \"id\"=\$1`).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectRollback()
			},
			id:    1,
			isErr: true,
		},
		{
			name: "Rows Affected error",
			mockBehavior: func(id int) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM \"assets\" WHERE \"id\"=\$1`).WillReturnResult(sqlmock.NewErrorResult(err))
				mock.ExpectRollback()
			},
			id:    1,
			isErr: true,
		},
	}
	for _, table := range testTables {
		t.Run(table.name, func(t *testing.T) {
			table.mockBehavior(table.id)

			err := repo.DeleteProduct(table.id)
			if table.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	logger := logger.GetLogger()

	product := &ProductRepo{
		db:     db,
		logger: logger,
		dbname: "assets",
	}

	repo := Repository{
		Product: product,
	}

	type mockBehavior func(product domain.Product)

	testTables := []struct {
		name         string
		mockBehavior mockBehavior
		product      domain.Product
		isErr        bool
	}{
		{
			name: "Valid",
			mockBehavior: func(product domain.Product) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE \"assets\" SET \"name\"=\$1, \"price\"=\$2 WHERE \"id\"=\$3`).
					WithArgs(product.Name, product.Price, product.Id).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			product: domain.Product{
				Id:    1,
				Name:  "name",
				Price: decimal.RequireFromString("50.00"),
			},
			isErr: false,
		},
		{
			name: "Transaction Error",
			mockBehavior: func(product domain.Product) {
				mock.ExpectBegin().WillReturnError(err)
				mock.ExpectRollback()
			},
			isErr: true,
		},
		{
			name: "Exec UPDATE Error",
			mockBehavior: func(product domain.Product) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE \"assets\" SET \"name\"=\$1, \"price\"=\$2 WHERE \"id\"=\$3`).
					WithArgs(product.Name, product.Price, product.Id).WillReturnError(err)
				mock.ExpectRollback()
			},
			product: domain.Product{
				Id:    1,
				Name:  "name",
				Price: decimal.RequireFromString("50.00"),
			},
			isErr: true,
		},
		{
			name: "Rows Affected Error",
			mockBehavior: func(product domain.Product) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE \"assets\" SET \"name\"=\$1, \"price\"=\$2 WHERE \"id\"=\$3`).
					WithArgs(product.Name, product.Price, product.Id).WillReturnResult(sqlmock.NewErrorResult(err))
				mock.ExpectRollback()
			},
			product: domain.Product{
				Id:    1,
				Name:  "name",
				Price: decimal.RequireFromString("50.00"),
			},
			isErr: true,
		},
	}
	for _, table := range testTables {
		t.Run(table.name, func(t *testing.T) {
			table.mockBehavior(table.product)

			err := repo.UpdateProduct(table.product)
			if table.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetAllProducts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	logger := logger.GetLogger()

	product := &ProductRepo{
		db:     db,
		logger: logger,
		dbname: "assets",
	}

	repo := Repository{
		Product: product,
	}

	type mockBehavior func(products []domain.Product)

	testTables := []struct {
		name         string
		mockBehavior mockBehavior
		products     []domain.Product
		isErr        bool
	}{
		{
			name: "Valid",
			mockBehavior: func(products []domain.Product) {
				mock.ExpectBegin()
				rows := mock.NewRows([]string{"id", "name", "price"}).
					AddRow(1, "name", "50.20").
					AddRow(2, "name2", "60.31").
					AddRow(3, "name3", "70.42")
				mock.ExpectQuery(`SELECT \* FROM \"assets\"`).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			products: []domain.Product{
				{
					Id:    1,
					Name:  "name",
					Price: decimal.RequireFromString("50.20"),
				},
				{
					Id:    2,
					Name:  "name2",
					Price: decimal.RequireFromString("60.31"),
				},
				{
					Id:    3,
					Name:  "name3",
					Price: decimal.RequireFromString("70.42"),
				},
			},
		},
		{
			name: "Transaction error",
			mockBehavior: func(products []domain.Product) {
				mock.ExpectBegin().WillReturnError(errors.New("transaction error"))
				mock.ExpectRollback()
			},
			isErr: true,
		},
		{
			name: "Query error",
			mockBehavior: func(products []domain.Product) {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT \* FROM \"assets\"`).WillReturnError(errors.New("query error"))
				mock.ExpectRollback()
			},
			isErr: true,
		},
		{
			name: "No Products",
			mockBehavior: func(products []domain.Product) {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT \* FROM \"assets\"`).WillReturnError(domain.ErrProductNotFound)
				mock.ExpectRollback()
			},
			isErr: true,
		},
	}
	for _, table := range testTables {
		t.Run(table.name, func(t *testing.T) {
			table.mockBehavior(table.products)

			got, err := repo.GetAllProducts()
			if table.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, table.products, got)
			}
		})
	}
}

func TestDeleteAllProducts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	logger := logger.GetLogger()

	product := &ProductRepo{
		db:     db,
		logger: logger,
		dbname: "assets",
	}

	repo := Repository{
		Product: product,
	}

	type mockBehavior func()

	testTables := []struct {
		name         string
		mockBehavior mockBehavior
		isErr        bool
	}{
		{
			name: "Valid",
			mockBehavior: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "assets"`).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			isErr: false,
		},
		{
			name: "Transaction error",
			mockBehavior: func() {
				mock.ExpectBegin().WillReturnError(err)
				mock.ExpectRollback()
			},
			isErr: true,
		},
		{
			name: "Exec DELETE error",
			mockBehavior: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "assets"`).WillReturnError(err)
				mock.ExpectRollback()
			},
			isErr: true,
		},
	}
	for _, table := range testTables {
		t.Run(table.name, func(t *testing.T) {
			table.mockBehavior()

			err := repo.DeleteAllProducts()
			if table.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
