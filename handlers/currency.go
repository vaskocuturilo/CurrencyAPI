package handlers

import (
	"CurrencyAPI/models"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
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

func (h *CurrencyHandler) GetCurrencies(c *gin.Context) {
	h.mu.Lock()
	defer h.mu.Unlock()

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, h.Cache)
}
