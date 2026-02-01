package service

import (
	"kasir-api/model"
	"kasir-api/repository"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll() ([]model.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoryService) Create(data *model.Category) error {
	return s.repo.Create(data)
}

func (s *CategoryService) GetByID(id int) (*model.Category, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) Update(data *model.Category) error {
	return s.repo.Update(data)
}

func (s *CategoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
