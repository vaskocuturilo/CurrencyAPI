package service

import (
	"CurrencyAPI/internal/models"
	"context"
	"errors"
	"testing"
)

type MockRepository struct {
	GetFunc     func(ctx context.Context, currency string) (*models.CurrencyResponse, error)
	GetAllFunc  func(ctx context.Context) ([]*models.CurrencyResponse, error)
	SaveAllFunc func(ctx context.Context, currencies []models.CurrencyResponse) error
}

func (m *MockRepository) Get(ctx context.Context, currency string) (*models.CurrencyResponse, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, currency)
	}
	return nil, nil
}

func (m *MockRepository) GetAll(ctx context.Context) ([]*models.CurrencyResponse, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return nil, nil
}

func (m *MockRepository) SaveAll(ctx context.Context, currencies []models.CurrencyResponse) error {
	if m.SaveAllFunc != nil {
		return m.SaveAllFunc(ctx, currencies)
	}
	return nil
}

func TestCurrencyService_Get_TableDriven(t *testing.T) {
	type testCase struct {
		name         string
		inputCode    string
		mockResponse *models.CurrencyResponse
		mockError    error
		wantErr      error
	}

	USD := models.CurrencyResponse{Date: "01-01-2026", Base: "USD", Quote: "EUR", Rate: 1.1}

	tests := []testCase{
		{
			name:         "Success - Valid 3-letter code",
			inputCode:    "USD",
			mockResponse: &USD,
			mockError:    nil,
			wantErr:      nil,
		},
		{
			name:         "Validation Error - Empty code",
			inputCode:    "",
			mockResponse: nil,
			mockError:    nil,
			wantErr:      models.ErrCurrencyRequired,
		},
		{
			name:         "Validation Error - Invalid length (2 chars)",
			inputCode:    "US",
			mockResponse: nil,
			mockError:    nil,
			wantErr:      models.ErrCurrencyInvalidLength,
		},
		{
			name:         "Repository Error - Not Found",
			inputCode:    "XYZ",
			mockResponse: nil,
			mockError:    models.ErrNotFound,
			wantErr:      models.ErrNotFound,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockRepository{
				GetFunc: func(ctx context.Context, currency string) (*models.CurrencyResponse, error) {
					return tc.mockResponse, tc.mockError
				},
			}

			serv := NewCurrencyService(mockRepo)
			// Act
			result, err := serv.Get(context.Background(), tc.inputCode)

			// Assert
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Expected error %v, but got %v", tc.wantErr, err)
			}

			if tc.wantErr == nil {
				if result == nil {
					t.Fatal("Expected result currency, but got nil")
				}
				if result.Base != tc.mockResponse.Base {
					t.Errorf("Expected Base %s, got %s", tc.mockResponse.Base, result.Base)
				}
			}
		})
	}
}

func TestCurrencyService_GetAll_TableDriven(t *testing.T) {
	type testCase struct {
		name         string
		mockResponse []*models.CurrencyResponse
		mockError    error
		expectedLen  int
		wantErr      error
	}

	currencies := []*models.CurrencyResponse{
		{Date: "2026-11-06", Base: "USD", Quote: "EUR", Rate: 1.1},
		{Date: "2026-11-06", Base: "EUR", Quote: "GBP", Rate: 0.85},
	}

	tests := []testCase{
		{
			name:         "Success - Returns all currencies",
			mockResponse: currencies,
			mockError:    nil,
			expectedLen:  2,
			wantErr:      nil,
		},
		{
			name:         "Success - Empty list",
			mockResponse: []*models.CurrencyResponse{},
			mockError:    nil,
			expectedLen:  0,
			wantErr:      nil,
		},
		{
			name:         "Error - Repository failure",
			mockResponse: nil,
			mockError:    errors.New("database connection failed"),
			expectedLen:  0,
			wantErr:      errors.New("database connection failed"),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockRepository{
				GetAllFunc: func(ctx context.Context) ([]*models.CurrencyResponse, error) {
					return tc.mockResponse, tc.mockError
				},
			}

			serv := NewCurrencyService(mockRepo)

			// Act
			result, err := serv.GetAll(context.Background())

			// Assert Error
			if (err != nil) != (tc.wantErr != nil) || (err != nil && err.Error() != tc.wantErr.Error()) {
				t.Fatalf("Expected error %v, got %v", tc.wantErr, err)
			}

			// Assert Length
			if len(result) != tc.expectedLen {
				t.Errorf("Expected %d currencies, got %d", tc.expectedLen, len(result))
			}
		})
	}
}

func TestCurrencyService_SaveAll_TableDriven(t *testing.T) {
	type testCase struct {
		name      string
		inputData []models.CurrencyResponse
		mockError error
		wantErr   error
	}

	sampleData := []models.CurrencyResponse{
		{Date: "2026-11-06", Base: "USD", Quote: "EUR", Rate: 1.1},
	}

	tests := []testCase{
		{
			name:      "Success - Saved successfully",
			inputData: sampleData,
			mockError: nil,
			wantErr:   nil,
		},
		{
			name:      "Error - Repository save failure",
			inputData: sampleData,
			mockError: errors.New("constraint violation"),
			wantErr:   errors.New("constraint violation"),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockRepository{
				SaveAllFunc: func(ctx context.Context, currencies []models.CurrencyResponse) error {
					return tc.mockError
				},
			}

			serv := NewCurrencyService(mockRepo)

			// Act
			err := serv.SaveAll(context.Background(), tc.inputData)

			// Assert Error
			if (err != nil) != (tc.wantErr != nil) || (err != nil && err.Error() != tc.wantErr.Error()) {
				t.Fatalf("Expected error %v, got %v", tc.wantErr, err)
			}
		})
	}
}
