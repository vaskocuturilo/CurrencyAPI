package handlers

import (
	"CurrencyAPI/models"
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type CurrencyHandler struct {
	mu    sync.RWMutex
	Cache []models.CurrencyResponse
}

func NewCurrencyHandler(cache []models.CurrencyResponse) *CurrencyHandler {
	return &CurrencyHandler{Cache: cache}
}

func (h *CurrencyHandler) UpdateCache(newCache []models.CurrencyResponse) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Cache = newCache
}

func (h *CurrencyHandler) GetCurrencies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(h.Cache); err != nil {
		log.Println("Failed to encode JSON")
		return
	}
}
