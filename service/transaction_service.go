package service

import (
	"kasir-api/model"
	"kasir-api/repository"
)

type TransactionService struct {
	repo *repository.TransactionRepository
}

func NewTransactionService(repo *repository.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []model.CheckoutItem) (*model.Transaction, error) {
	return s.repo.CreateTransaction(items)
}

func (s *TransactionService) GetTransactionReport() (*model.TransactionReport, error) {
	return s.repo.GetTransactionReport()
}

func (s *TransactionService) GetTransactionReportByDate(startDate string, endDate string) (*model.TransactionReport, error) {
	return s.repo.GetTransactionReportByDate(startDate, endDate)
}
