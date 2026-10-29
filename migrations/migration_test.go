package migrations

import (
	"database/sql"
	"os"
	"testing"
)

func TestRunMigrations_Integration(t *testing.T) {
	dbURL := os.Getenv("TEST_DB_URL")
	if dbURL == "" {
		t.Skip("Skipping migration integration test: TEST_DB_URL not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("test database is not reachable: %v", err)
	}

	_, _ = db.Exec("DROP TABLE IF EXISTS schema_migrations;")

	err = RunMigrations(db)
	if err != nil {
		t.Fatalf("expected migrations to run successfully on first try, got %v", err)
	}

	err = RunMigrations(db)
	if err != nil {
		t.Fatalf("expected migrations to succeed on second run (up to date), got %v", err)
	}
}
