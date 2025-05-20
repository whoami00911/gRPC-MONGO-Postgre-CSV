package service

import (
	"webApp/domain"
	"webApp/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

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
