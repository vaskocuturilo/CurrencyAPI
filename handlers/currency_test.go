package handlers

import (
	"CurrencyAPI/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCurrencyHandler_UpdateCache(t *testing.T) {
	initialCache := []models.CurrencyResponse{
		{Date: "2026-10-05", Base: "EUR", Quote: "USD", Rate: 1.10},
	}
	handler := NewCurrencyHandler(initialCache)

	newCache := []models.CurrencyResponse{
		{Date: "2026-10-05", Base: "EUR", Quote: "GBP", Rate: 0.85},
	}

	handler.UpdateCache(newCache)

	req := httptest.NewRequest(http.MethodGet, "/currencies", nil)
	rec := httptest.NewRecorder()

	handler.GetCurrencies(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response []models.CurrencyResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if len(response) != 1 || response[0].Quote != "GBP" {
		t.Errorf("expected updated cache with GBP, got %v", response)
	}
}
