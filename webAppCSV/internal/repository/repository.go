package repository

import (
	"database/sql"
	"webApp/domain"
	"webApp/pkg/logger"
)

type Product interface {
	AddProduct(product domain.Product) error
	GetProduct(id int) (domain.Product, error)
	DeleteProduct(id int) error
	UpdateProduct(product domain.Product) error
	GetAllProducts() ([]domain.Product, error)
	DeleteAllProducts() error
}
type Repository struct {
	Product
}

func NewRepository(db *sql.DB, logger *logger.Logger) *Repository {
	return &Repository{
		Product: newProductRepo(db, logger),
	}
}
