package service

import (
	"webApp/domain"
	"webApp/internal/repository"
)

type Product interface {
	AddProduct(product domain.Product) error
	GetProduct(id int) (domain.Product, error)
	DeleteProduct(id int) error
	UpdateProduct(product domain.Product) error
	GetAllProducts() ([]domain.Product, error)
	DeleteAllProducts() error
}

type Service struct {
	Product
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Product: newProductRepo(repo),
	}
}
