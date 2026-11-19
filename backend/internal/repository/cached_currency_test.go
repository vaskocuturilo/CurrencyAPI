package repository

import (
	"CurrencyAPI/internal/models"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

type stubCurrencyRepository struct {
	getFunc     func(ctx context.Context, currency string) (*models.CurrencyResponse, error)
	getAllFunc  func(ctx context.Context) ([]*models.CurrencyResponse, error)
	saveAllFunc func(ctx context.Context, currencies []models.CurrencyResponse) error
}

func (s *stubCurrencyRepository) Get(ctx context.Context, currency string) (*models.CurrencyResponse, error) {
	if s.getFunc != nil {
		return s.getFunc(ctx, currency)
	}
	return nil, nil
}

func (s *stubCurrencyRepository) GetAll(ctx context.Context) ([]*models.CurrencyResponse, error) {
	if s.getAllFunc != nil {
		return s.getAllFunc(ctx)
	}
	return nil, nil
}

func (s *stubCurrencyRepository) SaveAll(ctx context.Context, currencies []models.CurrencyResponse) error {
	if s.saveAllFunc != nil {
		return s.saveAllFunc(ctx, currencies)
	}
	return nil
}

func setupTestRedis(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	t.Helper()

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	t.Cleanup(mr.Close)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	t.Cleanup(func() {
		_ = client.Close()
	})

	return mr, client
}

func TestCachedCurrencyRepository_Get_CacheMissThenHit(t *testing.T) {
	mr, rdb := setupTestRedis(t)
	ctx := context.Background()

	dbCalls := 0
	expected := &models.CurrencyResponse{
		Date:  "2026-11-19",
		Base:  "USD",
		Quote: "EUR",
		Rate:  0.92,
	}

	stubRepo := &stubCurrencyRepository{
		getFunc: func(ctx context.Context, currency string) (*models.CurrencyResponse, error) {
			dbCalls++
			return expected, nil
		},
	}

	cachedRepo := NewCachedCurrencyRepository(stubRepo, rdb, 10*time.Minute)

	got1, err := cachedRepo.Get(ctx, "eur")
	if err != nil {
		t.Fatalf("unexpected error on cache miss: %v", err)
	}
	if dbCalls != 1 {
		t.Fatalf("expected 1 DB call on cache miss, got %d", dbCalls)
	}
	if got1 == nil || got1.Rate != expected.Rate || got1.Quote != "EUR" {
		t.Fatalf("unexpected currency returned: %+v", got1)
	}

	if !mr.Exists("currency:EUR") {
		t.Fatal("expected key currency:EUR to be stored in Redis")
	}
	if mr.TTL("currency:EUR") <= 0 {
		t.Fatal("expected positive TTL on cached key")
	}

	got2, err := cachedRepo.Get(ctx, "EUR")
	if err != nil {
		t.Fatalf("unexpected error on cache hit: %v", err)
	}
	if dbCalls != 1 {
		t.Fatalf("expected DB call count to remain 1 on cache hit, got %d", dbCalls)
	}
	if got2 == nil || got2.Rate != expected.Rate {
		t.Fatalf("expected cached rate %.2f, got %+v", expected.Rate, got2)
	}
}

func TestCachedCurrencyRepository_Get_DirectCacheHit(t *testing.T) {
	mr, rdb := setupTestRedis(t)
	ctx := context.Background()

	preloaded := models.CurrencyResponse{
		Date:  "2026-11-19",
		Base:  "USD",
		Quote: "PLN",
		Rate:  4.05,
	}
	raw, err := json.Marshal(preloaded)
	if err != nil {
		t.Fatalf("failed to marshal fixture: %v", err)
	}
	if err := mr.Set("currency:PLN", string(raw)); err != nil {
		t.Fatalf("failed to seed miniredis: %v", err)
	}

	dbCalled := false
	stubRepo := &stubCurrencyRepository{
		getFunc: func(ctx context.Context, currency string) (*models.CurrencyResponse, error) {
			dbCalled = true
			return nil, nil
		},
	}

	cachedRepo := NewCachedCurrencyRepository(stubRepo, rdb, 10*time.Minute)

	got, err := cachedRepo.Get(ctx, "PLN")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dbCalled {
		t.Fatal("expected repository not to be called on cache hit")
	}
	if got == nil || got.Rate != 4.05 || got.Quote != "PLN" {
		t.Fatalf("unexpected currency from cache: %+v", got)
	}
}

func TestCachedCurrencyRepository_SaveAll_SyncsRedis(t *testing.T) {
	_, rdb := setupTestRedis(t)
	ctx := context.Background()

	saved := false
	stubRepo := &stubCurrencyRepository{
		saveAllFunc: func(ctx context.Context, currencies []models.CurrencyResponse) error {
			saved = true
			return nil
		},
		getFunc: func(ctx context.Context, currency string) (*models.CurrencyResponse, error) {
			t.Fatal("expected Get to read synced value from Redis without calling DB")
			return nil, nil
		},
	}

	cachedRepo := NewCachedCurrencyRepository(stubRepo, rdb, 10*time.Minute)

	rates := []models.CurrencyResponse{
		{Date: "2026-11-19", Base: "USD", Quote: "EUR", Rate: 0.95},
		{Date: "2026-11-19", Base: "USD", Quote: "GBP", Rate: 0.79},
	}

	if err := cachedRepo.SaveAll(ctx, rates); err != nil {
		t.Fatalf("unexpected SaveAll error: %v", err)
	}
	if !saved {
		t.Fatal("expected underlying repository SaveAll to be called")
	}

	gotEUR, err := cachedRepo.Get(ctx, "EUR")
	if err != nil {
		t.Fatalf("unexpected Get error for EUR: %v", err)
	}
	if gotEUR == nil || gotEUR.Rate != 0.95 {
		t.Fatalf("expected synced EUR rate 0.95, got %+v", gotEUR)
	}
}
