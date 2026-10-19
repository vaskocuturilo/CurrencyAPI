package providers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchBNRRates(t *testing.T) {
	mockXML := `<?xml version="1.0" encoding="utf-8"?>
<DataSet xmlns="https://www.bnr.ro/xsd" xmlns:xsi="https://www.w3.org/2001/XMLSchema-instance">
    <Header>
        <Publisher>National Bank of Romania</Publisher>
        <PublishingDate>2026-10-19</PublishingDate>
        <MessageType>DR</MessageType>
    </Header>
    <Body>
        <Subject>Reference rates</Subject>
        <OrigCurrency>RON</OrigCurrency>
        <Cube date="2026-10-19">
            <Rate currency="USD">4.7520</Rate>
            <Rate currency="EUR">5.3272</Rate>
            <Rate currency="HUF" multiplier="100">1.4677</Rate>
        </Cube>
    </Body>
</DataSet>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockXML))
	}))
	defer server.Close()

	oldURL := bancaNationalRomaniaURL
	bancaNationalRomaniaURL = server.URL

	defer func() { bancaNationalRomaniaURL = oldURL }()

	rates, err := FetchBancaNationalRomaniaRates()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(rates) != 3 {
		t.Fatalf("expected 3 rates, got %d", len(rates))
	}

	// Validate USD mapping
	if rates[0].Quote != "USD" || rates[0].Rate != 4.7520 || rates[0].Base != "RON" || rates[0].Date != "2026-10-19" {
		t.Errorf("unexpected USD rate parsing: %+v", rates[0])
	}

	// Validate multi-unit normalization for HUF (1.4677 / 100 = 0.014677)
	expectedHUFRate := 1.4677 / 100.0
	if rates[2].Quote != "HUF" || rates[2].Rate != expectedHUFRate {
		t.Errorf("expected normalized HUF rate %f, got %f", expectedHUFRate, rates[2].Rate)
	}
}
