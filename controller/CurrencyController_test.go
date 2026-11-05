package controller

import (
	"CurrencyAPI/models"
	"context"
	"testing"
)

type MockCurrencyService struct {
	GetFunc     func(ctx context.Context, currency string) (*models.CurrencyResponse, error)
	GetAllFunc  func(ctx context.Context) ([]*models.CurrencyResponse, error)
	SaveAllFunc func(ctx context.Context, currencies []models.CurrencyResponse) error
}

func (m *MockCurrencyService) Get(ctx context.Context, currency string) (*models.CurrencyResponse, error) {
	if m.GetFunc != nil {
		return m.Get(ctx, currency)
	}
	return nil, nil
}

func (m *MockCurrencyService) GetAll(ctx context.Context) ([]*models.CurrencyResponse, error) {
	if m.GetAllFunc != nil {
		return m.GetAll(ctx)
	}
	return nil, nil
}

func (m *MockCurrencyService) SaveAll(ctx context.Context, currencies []models.CurrencyResponse) error {
	if m.SaveAllFunc != nil {
		return m.SaveAll(ctx, currencies)
	}
	return nil
}

func TestCurrencyController_Get_Table(t *testing.T)        {}
func TestCurrencyController_GetAll_Table(t *testing.T)     {}
func TestCurrencyController_GetSaveAll_Table(t *testing.T) {}
