package handlers

import (
	"CurrencyAPI/models"
	"log/slog"
	"net/http"
	"sync"
	"time"

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
	start := time.Now()

	h.mu.RLock()
	cacheLen := len(h.Cache)
	h.mu.RUnlock()

	if cacheLen == 0 {
		slog.Error("The currency cache is empty during request", "path",
			c.FullPath(), "client_ip", c.ClientIP())
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Service is initializing or cache is empty"})
		return
	}

	h.mu.RLock()
	data := h.Cache
	h.mu.RUnlock()
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, data)

	slog.Info("Handled get currencies request",
		"method", c.Request.Method,
		"path", c.FullPath(),
		"status", http.StatusOK,
		"items_returned", cacheLen,
		"latency_ms", time.Since(start).Milliseconds(),
		"client_ip", c.ClientIP(),
	)
}
