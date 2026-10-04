package handlers

import (
	"CurrencyAPI/models"
	"encoding/json"
	"net/http"
)

type CurrencyHandler struct {
	Cache []models.CurrencyResponse
}

func NewCurrencyHandler(cache []models.CurrencyResponse) *CurrencyHandler {
	return &CurrencyHandler{Cache: cache}
}

func (h *CurrencyHandler) GetCurrencies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(h.Cache); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
