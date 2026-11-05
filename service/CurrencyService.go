package service

import (
	"CurrencyAPI/models"
	"CurrencyAPI/repository"
	"context"
)

type ICurrencyService interface {
	Get(ctx context.Context, currency string) (*models.CurrencyResponse, error)
	GetAll(ctx context.Context) ([]*models.CurrencyResponse, error)
	SaveAll(ctx context.Context, currencies []models.CurrencyResponse) error
}

type CurrencyService struct {
	repo repository.PostgresCurrencyRepository
}

func NewCurrencyService(repo repository.PostgresCurrencyRepository) *CurrencyService {
	return &CurrencyService{repo: repo}
}

func (c *CurrencyService) Get(ctx context.Context, currency string) (*models.CurrencyResponse, error) {
	return nil, nil
}

func (c *CurrencyService) GetAll(ctx context.Context) ([]*models.CurrencyResponse, error) {
	return nil, nil
}

func (c *CurrencyService) SaveAll(ctx context.Context, currencies []models.CurrencyResponse) error {
	return nil
}
