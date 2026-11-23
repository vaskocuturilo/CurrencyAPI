package providers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type rewriteTransport struct {
	targetBase string
	rt         http.RoundTripper
}

func (t *rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	clone := req.Clone(req.Context())
	idx := strings.Index(clone.URL.String(), "/api/v2/")
	if idx != -1 {
		suffix := clone.URL.String()[idx:]
		clone.URL = nil
		var err error
		req = clone
		_ = err
		parsed, err := req.URL.Parse(t.targetBase + suffix)
		if err != nil {
			return nil, err
		}
		req.URL = parsed
	}
	if t.rt != nil {
		return t.rt.RoundTrip(req)
	}
	return http.DefaultTransport.RoundTrip(req)
}

func swapHTTPClientForTest(t *testing.T, serverURL string) {
	t.Helper()
	origClient := httpClient
	httpClient = &http.Client{
		Transport: &rewriteTransport{
			targetBase: serverURL,
			rt:         http.DefaultTransport,
		},
	}
	t.Cleanup(func() {
		httpClient = origClient
	})
}

func TestFetchBankInternationalSettlements_SuccessWithCurrencyCol(t *testing.T) {
	var seenAccept string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenAccept = r.Header.Get("Accept")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(
			"FREQ,CURRENCY,TIME_PERIOD,OBS_VALUE\n" +
				"D,eur,2026-11-20,0.92\n" +
				"D,GBP,2026-11-20,0.79\n" +
				"D,INVALID_CODE,2026-11-20,1.50\n" +
				"D,JPY,2026-11-20,not_a_number\n",
		))
	}))
	defer srv.Close()

	swapHTTPClientForTest(t, srv.URL)

	rates, err := FetchBankInternationalSettlements()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if seenAccept != "application/vnd.sdmx.data+csv;version=1.0.0" {
		t.Errorf("unexpected Accept header: %q", seenAccept)
	}
	if len(rates) != 2 {
		t.Fatalf("expected 2 valid rates, got %d (%+v)", len(rates), rates)
	}
	if rates[0].Quote != "EUR" || rates[0].Base != "USD" || rates[0].Rate != 0.92 || rates[0].Date != "2026-11-20" {
		t.Errorf("unexpected first rate: %+v", rates[0])
	}
	if rates[1].Quote != "GBP" || rates[1].Rate != 0.79 {
		t.Errorf("unexpected second rate: %+v", rates[1])
	}
}

func TestFetchBankInternationalSettlements_RefAreaFallbackAndMalformedRow(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(
			"REF_AREA,TIME_PERIOD,OBS_VALUE\n" +
				"CHF,2026-11-20,0.88\n" +
				"\"unclosed_quote,2026-11-20,1.00\n",
		))
	}))
	defer srv.Close()

	swapHTTPClientForTest(t, srv.URL)

	rates, err := FetchBankInternationalSettlements()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rates) != 1 || rates[0].Quote != "CHF" || rates[0].Rate != 0.88 {
		t.Fatalf("expected 1 CHF rate, got %+v", rates)
	}
}

func TestFetchBankInternationalSettlements_Errors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
	}{
		{
			name:       "non-200 status",
			statusCode: http.StatusInternalServerError,
			body:       "upstream error",
		},
		{
			name:       "empty body",
			statusCode: http.StatusOK,
			body:       "",
		},
		{
			name:       "missing required columns",
			statusCode: http.StatusOK,
			body:       "FOO,BAR\n1,2\n",
		},
		{
			name:       "missing currency and ref_area columns",
			statusCode: http.StatusOK,
			body:       "TIME_PERIOD,OBS_VALUE\n2026-11-20,1.25\n",
		},
		{
			name:       "no valid rates in rows",
			statusCode: http.StatusOK,
			body:       "CURRENCY,TIME_PERIOD,OBS_VALUE\nINVALID,2026-11-20,1.25\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer srv.Close()

			swapHTTPClientForTest(t, srv.URL)

			rates, err := FetchBankInternationalSettlements()
			if err == nil {
				t.Fatalf("expected error for case %q, got rates: %+v", tc.name, rates)
			}
		})
	}
}
