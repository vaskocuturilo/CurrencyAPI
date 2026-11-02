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

	id := "123"

	rows := sqlmock.NewRows([]string{"base", "rate", "quote", "date"}).
		AddRow(id, "USD", 1.1)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, base, rate, quote, date FROM currencies_db WHERE id = $1")).
		WithArgs(id).
		WillReturnRows(rows)

	currency, err := repo.Get(context.Background(), id)

	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}
	if currency.Base != "USD" {
		t.Errorf("expected base 'USD', got %s", currency.Base)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresEventRepository_Get_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostgresCurrencyRepository(db)

	mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)

	_, err := repo.Get(context.Background(), "unknown")

	if !errors.Is(err, models.ErrNotFound) {
		t.Errorf("expected models.ErrNotFound, got %v", err)
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
		{Base: "USD", Rate: 1.1, Date: "Desc 1", Quote: "1"},
		{Base: "EUR", Rate: 1.1, Date: "Desc 2", Quote: "2"},
	}

	rows := sqlmock.NewRows([]string{"base", "rate", "date", "quote"})

	for _, c := range expectedCurrencies {
		rows.AddRow(c.Base, c.Quote, c.Rate, c.Date)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT *")).WillReturnRows(rows)

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
