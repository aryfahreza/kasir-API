package service

import (
	"kasir-api/model"
	"kasir-api/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll() ([]model.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) Create(data *model.Product) error {
	return s.repo.Create(data)
}

func (s *ProductService) GetByID(id int) (*model.Product, error) {
	return s.repo.GetByID(id)
}
