package repository

import (
	"CurrencyAPI/models"
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPostgresEventRepository_Get_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresCurrencyRepository(db)

	currency := "USD"

	rows := sqlmock.NewRows([]string{"date", "base", "quote", "rate"}).
		AddRow(currency, "2026-11-03", "USD", 1.1)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT date, base, quote,rate FROM currencies WHERE quote = $1")).
		WithArgs(currency).
		WillReturnRows(rows)

	currencyResponse, err := repo.Get(context.Background(), currency)

	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}
	if currencyResponse.Quote != "USD" {
		t.Errorf("expected base 'USD', got %s", currencyResponse.Quote)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresEventRepository_Get_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostgresCurrencyRepository(db)

	mock.ExpectQuery("SELECT *").WillReturnError(sql.ErrNoRows)

	_, err := repo.Get(context.Background(), "unknown")

	if !errors.Is(err, models.ErrNotFound) {
		t.Errorf("expected domain.ErrNotFound, got %v", err)
	}
}

func TestPostgresEventRepository_GetAll_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresCurrencyRepository(db)

	expectedCurrencies := []models.CurrencyResponse{
		{Date: "2026-11-03", Base: "USD", Quote: "EUR", Rate: 1.1},
		{Date: "2026-11-03", Base: "EUR", Quote: "GBP", Rate: 0.85},
	}

	rows := sqlmock.NewRows([]string{"date", "base", "quote", "rate"})

	for _, c := range expectedCurrencies {
		rows.AddRow(c.Date, c.Base, c.Quote, c.Rate)
	}

	queryRegex := regexp.QuoteMeta(`SELECT date, base, quote, rate FROM currencies`)
	mock.ExpectQuery(queryRegex).WillReturnRows(rows)

	currencies, err := repo.GetAll(context.Background())

	if len(currencies) != len(expectedCurrencies) {
		t.Fatalf("expected %d currencies, got %d", len(expectedCurrencies), len(currencies))
	}

	for i, e := range expectedCurrencies {
		if currencies[i].Base != e.Base {
			t.Errorf("[%d] expected Base %s, got %s", i, e.Base, currencies[i].Base)
		}

		if currencies[i].Quote != e.Quote {
			t.Errorf("[%d] expected Quote %s, got %s", i, e.Quote, currencies[i].Quote)
		}

		if currencies[i].Rate != e.Rate {
			t.Errorf("[%d] expected Rate %f, got %f", i, e.Rate, currencies[i].Rate)
		}

		if currencies[i].Date != e.Date {
			t.Errorf("[%d] expected Date %s, got %s", i, e.Date, currencies[i].Date)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %s", err)
	}
}

func TestPostgresEventRepository_GetAll_Empty(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostgresCurrencyRepository(db)

	rows := sqlmock.NewRows([]string{"id", "title", "description"})

	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	currency, err := repo.GetAll(context.Background())

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(currency) != 0 {
		t.Errorf("expected 0 events, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %s", err)
	}
}

func TestPostgresEventRepository_Save_All(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresCurrencyRepository(db)

	expectedCurrencies := []models.CurrencyResponse{
		{Base: "USD", Rate: 1.1, Date: "Desc 1", Quote: "1"},
		{Base: "EUR", Rate: 1.1, Date: "Desc 2", Quote: "2"},
	}

	queryRegex := regexp.QuoteMeta(`
		INSERT INTO currencies (date, base, quote, rate)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (date, base, quote) 
		DO UPDATE SET rate = EXCLUDED.rate
	`)

	mock.ExpectPrepare(queryRegex)

	mock.ExpectExec(queryRegex).
		WithArgs(expectedCurrencies[0].Date, expectedCurrencies[0].Base, expectedCurrencies[0].Quote, expectedCurrencies[0].Rate).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(queryRegex).
		WithArgs(expectedCurrencies[1].Date, expectedCurrencies[1].Base, expectedCurrencies[1].Quote, expectedCurrencies[1].Rate).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.SaveAll(context.Background(), expectedCurrencies)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %s", err)
	}
}
