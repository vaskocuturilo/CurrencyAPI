package repository

import (
	"CurrencyAPI/models"
	"context"
)

type CurrencyRepository interface {
	Get(ctx context.Context, currency string) (*models.CurrencyResponse, error)
	GetAll(ctx context.Context) ([]*models.CurrencyResponse, error)
}
