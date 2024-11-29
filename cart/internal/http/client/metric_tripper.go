package client

import (
	"net/http"
	"route256/cart/internal/metrics"
	"time"
)

type MetricTripper struct {
	transport http.RoundTripper
}

func NewMetricTripper(transport http.RoundTripper) http.RoundTripper {
	return &MetricTripper{
		transport: transport,
	}
}

func (m MetricTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	// Засекаем время начала запроса
	start := time.Now()

	resp, err := m.transport.RoundTrip(request)

	// Длительность запроса
	duration := time.Since(start)
	if err != nil {
		metrics.RecordExternalRequest(request.URL.Path, "error", duration)
		return nil, err
	}

	metrics.RecordExternalRequest(request.URL.Path, resp.Status, duration)

	return resp, err
}
