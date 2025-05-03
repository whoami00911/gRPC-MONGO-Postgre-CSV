package service

import (
	"webApp/domain"
	"webApp/internal/repository"
)

type ProductRepository struct {
	repo repository.Product
}

func newProductRepo(repo repository.Product) *ProductRepository {
	return &ProductRepository{
		repo: repo,
	}
}

func (p *ProductRepository) AddProduct(product domain.Product) error {
	return p.repo.AddProduct(product)
}

func (p *ProductRepository) GetProduct(id int) (domain.Product, error) {
	return p.repo.GetProduct(id)
}

func (p *ProductRepository) DeleteProduct(id int) error {
	return p.repo.DeleteProduct(id)
}

func (p *ProductRepository) UpdateProduct(product domain.Product) error {
	return p.repo.UpdateProduct(product)
}

func (p *ProductRepository) GetAllProducts() ([]domain.Product, error) {
	return p.repo.GetAllProducts()
}

func (p *ProductRepository) DeleteAllProducts() error {
	return p.repo.DeleteAllProducts()
}
