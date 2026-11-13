package providers

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func fetchURL(url string, headers map[string]string) ([]byte, error) {
	start := time.Now()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		slog.Error("Failed to create HTTP request for provider", "url", url, "error", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, val := range headers {
		req.Header.Set(key, val)
	}

	resp, err := httpClient.Do(req)

	if err != nil {
		slog.Error("Network error during external provider fetch", "url", url, "error", err)
		return nil, fmt.Errorf("network error: %w", err)
	}

	defer func(Body io.ReadCloser) {
		if closeErr := Body.Close(); closeErr != nil {
			slog.Warn("Failed to close response body", "url", url, "error", closeErr)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		slog.Error("External API returned non-200 status",
			"url", url,
			"status_code", resp.StatusCode,
			"status", resp.Status,
			"latency_ms", time.Since(start).Milliseconds(),
		)
		return nil, fmt.Errorf("API returned non-200 status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read response body from external provider")
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	slog.Info("Successfully fetched data from external provider URL",
		"url", url,
		"status", resp.Status,
		"bytes_read", len(body),
		"latency_ms", time.Since(start).Milliseconds(),
	)

	return body, nil
}
