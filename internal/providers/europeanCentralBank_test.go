package providers

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func helperTestClient(fn roundTripFunc) {
	httpClient = &http.Client{
		Transport: fn,
	}
}

func TestFetchEuropeanCentraBankRates_Success(t *testing.T) {
	mockXML := `<?xml version="1.0" encoding="UTF-8"?>
	<Envelope>
		<Cube>
			<Cube time="2026-10-28">
				<Cube currency="USD" rate="1.0850"/>
				<Cube currency="GBP" rate="0.8520"/>
			</Cube>
		</Cube>
	</Envelope>`

	helperTestClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(mockXML))),
		}, nil
	})

	currencies, err := FetchEuropeanCentraBankRates()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(currencies) != 2 {
		t.Fatalf("expected 2 currency rates, got %d", len(currencies))
	}

	if currencies[0].Quote != "USD" || currencies[0].Rate != 1.085 || currencies[0].Date != "2026-10-28" {
		t.Errorf("unexpected parsed currency data for USD: %+v", currencies[0])
	}

	if currencies[1].Quote != "GBP" || currencies[1].Rate != 0.852 || currencies[1].Date != "2026-10-28" {
		t.Errorf("unexpected parsed currency data for GBP: %+v", currencies[1])
	}
}

func TestFetchEuropeanCentraBankRates_Non200Status(t *testing.T) {
	helperTestClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Status:     "500 Internal Server Error",
			Body:       io.NopCloser(bytes.NewReader([]byte("Internal Error"))),
		}, nil
	})

	_, err := FetchEuropeanCentraBankRates()
	if err == nil {
		t.Fatal("expected an error due to non-200 status code, got nil")
	}
}

func TestFetchEuropeanCentraBankRates_InvalidXML(t *testing.T) {
	helperTestClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte("<InvalidXML>..."))),
		}, nil
	})

	_, err := FetchEuropeanCentraBankRates()
	if err == nil {
		t.Fatal("expected an error due to malformed XML, got nil")
	}
}
