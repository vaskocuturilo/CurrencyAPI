package controller

import (
	"CurrencyAPI/models"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type MockCurrencyService struct {
	GetFunc     func(ctx context.Context, currency string) (*models.CurrencyResponse, error)
	GetAllFunc  func(ctx context.Context) ([]*models.CurrencyResponse, error)
	SaveAllFunc func(ctx context.Context, currencies []models.CurrencyResponse) error
}

func (m *MockCurrencyService) Get(ctx context.Context, currency string) (*models.CurrencyResponse, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, currency)
	}
	return nil, nil
}

func (m *MockCurrencyService) GetAll(ctx context.Context) ([]*models.CurrencyResponse, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return nil, nil
}

func (m *MockCurrencyService) SaveAll(ctx context.Context, currencies []models.CurrencyResponse) error {
	if m.SaveAllFunc != nil {
		return m.SaveAllFunc(ctx, currencies)
	}
	return nil
}

func TestCurrencyController_Get_Table(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type testCase struct {
		name           string
		queryString    string
		mockResponse   *models.CurrencyResponse
		mockError      error
		expectedStatus int
	}

	sampleCurrency := &models.CurrencyResponse{Date: "2026-11-09", Base: "USD", Quote: "EUR", Rate: 1.1}

	tests := []testCase{
		{
			name:           "Success - Found currency by scope",
			queryString:    "?scope=EUR",
			mockResponse:   sampleCurrency,
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error - Currency not found",
			queryString:    "?scope=XYZ",
			mockResponse:   nil,
			mockError:      models.ErrNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Error - Internal server error",
			queryString:    "?scope=EUR",
			mockResponse:   nil,
			mockError:      errors.New("db error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockService := &MockCurrencyService{
				GetFunc: func(ctx context.Context, currency string) (*models.CurrencyResponse, error) {
					return tc.mockResponse, tc.mockError
				},
			}

			controller := NewUserController(mockService)

			router := gin.New()
			router.GET("/currencies", controller.Route)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/currencies"+tc.queryString, nil)
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d (body: %s)", tc.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestCurrencyController_GetAll_Table(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type testCase struct {
		name           string
		mockResponse   []*models.CurrencyResponse
		mockError      error
		expectedStatus int
	}

	sampleCurrencies := []*models.CurrencyResponse{
		{Date: "2026-11-09", Base: "USD", Quote: "EUR", Rate: 1.1},
	}

	tests := []testCase{
		{
			name:           "Success - Returns all currencies",
			mockResponse:   sampleCurrencies,
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error - Internal server error",
			mockResponse:   nil,
			mockError:      errors.New("db error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockService := &MockCurrencyService{
				GetAllFunc: func(ctx context.Context) ([]*models.CurrencyResponse, error) {
					return tc.mockResponse, tc.mockError
				},
			}

			controller := NewUserController(mockService)

			router := gin.New()
			router.GET("/currencies", controller.Route)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/currencies", nil)
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}
		})
	}
}
func TestCurrencyController_GetSaveAll_Table(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type testCase struct {
		name           string
		mockError      error
		expectedStatus int
	}

	tests := []testCase{
		{
			name:           "Success - Saved successfully",
			mockError:      nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Error - Save failure",
			mockError:      errors.New("save error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockService := &MockCurrencyService{
				SaveAllFunc: func(ctx context.Context, currencies []models.CurrencyResponse) error {
					return tc.mockError
				},
			}

			controller := NewUserController(mockService)

			router := gin.New()
			router.POST("/currencies", controller.SaveAll)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/currencies", nil)
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}
		})
	}
}
