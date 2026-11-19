package repository

import (
	"CurrencyAPI/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

const currencyKeyPrefix = "currency:%s"

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
	code := strings.ToUpper(strings.TrimSpace(currency))
	key := fmt.Sprintf(currencyKeyPrefix, code)

	val, err := c.cache.Get(ctx, key).Result()
	if err == nil {
		var cached models.CurrencyResponse
		if unmarshalErr := json.Unmarshal([]byte(val), &cached); unmarshalErr == nil {
			return &cached, nil
		}
		slog.Warn("Failed to unmarshal cached currency, falling back to DB", "key", key)
	} else if !errors.Is(err, redis.Nil) {
		slog.Warn("Redis GET error, falling back to DB", "key", key, "error", err)
	}

	res, err, _ := c.group.Do(key, func() (any, error) {
		dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		item, dbErr := c.repo.Get(dbCtx, code)
		if dbErr != nil {
			return nil, dbErr
		}
		if item == nil {
			return nil, nil
		}

		payload, marshalErr := json.Marshal(item)
		if marshalErr != nil {
			slog.Warn("Failed to marshal currency for cache", "key", key, "error", marshalErr)
			return item, nil
		}

		cacheCtx, cacheCancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cacheCancel()

		if setErr := c.cache.Set(cacheCtx, key, payload, c.ttl).Err(); setErr != nil {
			slog.Warn("Failed to write currency to Redis", "key", key, "error", setErr)
		}

		return item, nil
	})

	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}

	return res.(*models.CurrencyResponse), nil
}

func (c *CachedCurrencyRepository) GetAll(ctx context.Context) ([]*models.CurrencyResponse, error) {
	return c.repo.GetAll(ctx)
}

func (c *CachedCurrencyRepository) SaveAll(ctx context.Context, currencies []models.CurrencyResponse) error {
	if err := c.repo.SaveAll(ctx, currencies); err != nil {
		return err
	}

	for _, item := range currencies {
		code := strings.ToUpper(strings.TrimSpace(item.Quote))
		if code == "" {
			continue
		}
		key := fmt.Sprintf(currencyKeyPrefix, code)

		payload, err := json.Marshal(item)
		if err != nil {
			slog.Warn("Failed to marshal currency during cache sync", "key", key, "error", err)
			continue
		}

		if err := c.cache.Set(ctx, key, payload, c.ttl).Err(); err != nil {
			slog.Warn("Failed to update Redis key after SaveAll", "key", key, "error", err)
		}
	}

	return nil
}
