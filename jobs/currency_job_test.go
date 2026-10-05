package jobs

import (
	"CurrencyAPI/handlers"
	"testing"
)

func TestCurrencyUpdateJob_Contract(t *testing.T) {
	handler := handlers.NewCurrencyHandler(nil)
	job := NewCurrencyUpdateJob(handler)

	expectedName := "Currency Update Job"
	if job.Name() != expectedName {
		t.Errorf("expected job name '%s', got '%s'", expectedName, job.Name())
	}

	if job.Schedule() == "" {
		t.Errorf("expected a valid cron schedule expression, got empty string")
	}
}
