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
	repo repository.CurrencyRepository
}

func NewCurrencyService(repo repository.CurrencyRepository) *CurrencyService {
	return &CurrencyService{repo: repo}
}

func (c *CurrencyService) Get(ctx context.Context, currency string) (*models.CurrencyResponse, error) {
	if err := models.Validate(currency); err != nil {
		return nil, err
	}

	return c.repo.Get(ctx, currency)
}

func (c *CurrencyService) GetAll(ctx context.Context) ([]*models.CurrencyResponse, error) {
	return c.repo.GetAll(ctx)
}

func (c *CurrencyService) SaveAll(ctx context.Context, currencies []models.CurrencyResponse) error {
	return c.repo.SaveAll(ctx, currencies)
}
