package service

import (
	"CurrencyAPI/models"
	"context"
	"testing"
)

type MockRepository struct {
	GetFunc     func(ctx context.Context, currency string) (*models.CurrencyResponse, error)
	GetAllFunc  func(ctx context.Context) ([]*models.CurrencyResponse, error)
	SaveAllFunc func(ctx context.Context, currencies []models.CurrencyResponse) error
}

func (m *MockRepository) Get(ctx context.Context, currency string) (*models.CurrencyResponse, error) {
	if m.GetFunc != nil {
		return m.Get(ctx, currency)
	}
	return nil, nil
}

func (m *MockRepository) GetAll(ctx context.Context) ([]*models.CurrencyResponse, error) {
	if m.GetAllFunc != nil {
		return m.GetAll(ctx)
	}
	return nil, nil
}

func (m *MockRepository) SaveAll(ctx context.Context, currencies []models.CurrencyResponse) error {
	if m.SaveAllFunc != nil {
		return m.SaveAll(ctx, currencies)
	}
	return nil
}

func TestCurrencyService_Get_TableDriven(t *testing.T)     {}
func TestCurrencyService_GetAll_TableDriven(t *testing.T)  {}
func TestCurrencyService_SaveAll_TableDriven(t *testing.T) {}
