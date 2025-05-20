package repository

import (
	"database/sql"
	"webApp/domain"
	"webApp/pkg/logger"

	"github.com/spf13/viper"
)

type ProductRepo struct {
	db     *sql.DB
	logger *logger.Logger
	dbname string
}

func newProductRepo(db *sql.DB, logger *logger.Logger) *ProductRepo {
	return &ProductRepo{
		db:     db,
		logger: logger,
		dbname: viper.GetString("db.name"),
	}
}

func (p *ProductRepo) AddProduct(product domain.Product) error {
	tx, err := p.db.Begin()
	if err != nil {
		p.logger.Errorf("Transaction not started: %s", err)
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			p.logger.Errorf("Something wrong with transaction: %s", err)
		} else {
			tx.Commit()
		}
	}()
	result, err := tx.Exec(`INSERT INTO "`+p.dbname+`" (name, price) VALUES ($1, $2)`, product.Name, product.Price)
	if err != nil {
		p.logger.Errorf("Add product error: %s", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		p.logger.Errorf("Can't get rows affected parameter: %s", err)
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrProductExists
	}
	return err
}

func (p *ProductRepo) GetProduct(id int) (domain.Product, error) {
	var product domain.Product
	tx, err := p.db.Begin()
	if err != nil {
		p.logger.Errorf("Transaction not started: %s", err)
		return domain.Product{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			p.logger.Errorf("Something wrong with trunsaction: %s", err)
		} else {
			tx.Commit()
		}
	}()
	query := tx.QueryRow(`SELECT "id", "name", "price" WHERE "id" = $1`, id)
	err = query.Scan(&product.Id, &product.Name, &product.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Product{}, domain.ErrProductNotFound
		}
		p.logger.Errorf("Can't get product from DB: %s", err)
		return domain.Product{}, err
	}
	return product, err
}

func (p *ProductRepo) DeleteProduct(id int) error {
	tx, err := p.db.Begin()
	if err != nil {
		p.logger.Errorf("Transaction not started: %s", err)
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			p.logger.Errorf("Something wrong with trunsaction: %s", err)
		} else {
			tx.Commit()
		}
	}()
	result, err := tx.Exec(`DELETE FROM "`+p.dbname+`" WHERE "id"=$1`, id)
	if err != nil {
		p.logger.Errorf("Delete product error: %s", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		p.logger.Errorf("Can't get rows affected parameter: %s", err)
	}
	if rowsAffected == 0 {
		return domain.ErrProductNotFound
	}
	return err
}

func (p *ProductRepo) UpdateProduct(product domain.Product) error {
	tx, err := p.db.Begin()
	if err != nil {
		p.logger.Errorf("Transaction not started: %s", err)
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			p.logger.Errorf("Something wrong with trunsaction: %s", err)
		} else {
			tx.Commit()
		}
	}()
	result, err := tx.Exec(`UPDATE "`+p.dbname+`" SET "name"=$1, "price"=$2 WHERE "id"=$3`, product.Name, product.Price, product.Id)
	if err != nil {
		p.logger.Errorf("Delete product error: %s", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		p.logger.Errorf("Can't get rows affected parameter: %s", err)
	}
	if rowsAffected == 0 {
		return domain.ErrProductNotFound
	}
	return err
}

func (p *ProductRepo) GetAllProducts() ([]domain.Product, error) {
	var products []domain.Product
	tx, err := p.db.Begin()
	if err != nil {
		p.logger.Errorf("Transaction not started: %s", err)
		return []domain.Product{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			p.logger.Errorf("Something wrong with trunsaction: %s", err)
		} else {
			tx.Commit()
		}
	}()

	query, err := tx.Query(`SELECT * FROM "` + p.dbname + `"`)
	if err != nil {
		p.logger.Errorf("Select all error: %s", err)
		return []domain.Product{}, err
	}

	for query.Next() {
		var product domain.Product
		err = query.Scan(&product.Id, &product.Name, &product.Price)
		if err != nil {
			if err == sql.ErrNoRows {
				return []domain.Product{}, domain.ErrProductNotFound
			}
			p.logger.Errorf("Can't get product from DB: %s", err)
			return []domain.Product{}, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (p *ProductRepo) DeleteAllProducts() error {
	tx, err := p.db.Begin()
	if err != nil {
		p.logger.Errorf("Transaction not started: %s", err)
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			p.logger.Errorf("Something wrong with trunsaction: %s", err)
		} else {
			tx.Commit()
		}
	}()
	_, err = tx.Exec(`DELETE FROM "` + p.dbname + `"`)
	if err != nil {
		p.logger.Errorf("Delete products error: %s", err)
		return err
	}
	return err
}
