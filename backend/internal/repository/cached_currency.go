package repository

import (
	"CurrencyAPI/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

const currencyString = "currency:%s"

type CachedCurrencyRepository struct {
	repo  CurrencyRepository
	cache *redis.Client
	group singleflight.Group
	ttl   time.Duration
}

func NewCachedCurrencyRepository(repo CurrencyRepository, cache *redis.Client, ttl time.Duration) *CachedCurrencyRepository {
	return &CachedCurrencyRepository{
		repo:  repo,
		cache: cache,
		ttl:   ttl,
	}
}

func (c *CachedCurrencyRepository) Get(ctx context.Context, currency string) (*models.CurrencyResponse, error) {
	key := fmt.Sprintf(currencyString, currency)

	val, err := c.cache.Get(ctx, key).Result()

	if err == nil {
		var currency models.CurrencyResponse
		if err := json.Unmarshal([]byte(val), &currency); err == nil {
			return &currency, nil
		}
	}

	if err != nil && !errors.Is(err, redis.Nil) {
		slog.Info("Redis error for key ", "key", "error", key, err)
	}

	result, err, _ := c.group.Do(currency, func() (interface{}, error) {
		dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		currency, err := c.repo.Get(dbCtx, currency)
		if err != nil {
			return nil, err
		}

		data, _ := json.Marshal(currency)
		c.cache.Set(context.Background(), key, data, c.ttl)

		return currency, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*models.CurrencyResponse), nil
}

func (c *CachedCurrencyRepository) GetAll(ctx context.Context) ([]*models.CurrencyResponse, error) {
	return c.repo.GetAll(ctx)
}
