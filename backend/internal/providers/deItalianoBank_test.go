package providers

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

func TestFetchBancaDItaliaRates_Success(t *testing.T) {
	mockCSV := `# Some comment line
Date,Col1,ISO,Col3,Rate,Col5,Col6,Col7,RefDate
2026-10-28,x,USD,x,1.0850,x,x,x,2026-10-28
2026-10-28,x,GBP,x,0.8520,x,x,x,2026-10-28
2026-10-28,x,INVALID,x,N.A.,x,x,x,2026-10-28
`

	helperTestClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(mockCSV))),
		}, nil
	})

	currencies, err := FetchBancaDItaliaRates()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(currencies) != 2 {
		t.Fatalf("expected 2 valid currencies, got %d", len(currencies))
	}

	if currencies[0].Quote != "USD" || currencies[0].Rate != 1.085 || currencies[0].Base != "EUR" {
		t.Errorf("unexpected USD parsed data: %+v", currencies[0])
	}
}

func TestFetchBancaDItaliaRates_InvalidCSV(t *testing.T) {
	helperTestClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Status:     "500 Internal Server Error",
			Body:       io.NopCloser(bytes.NewReader([]byte("Error"))),
		}, nil
	})

	_, err := FetchBancaDItaliaRates()
	if err == nil {
		t.Fatal("expected an error due to non-200 status code, got nil")
	}
}
