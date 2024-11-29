package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

// Инициализация счетчиков, гистограмм и меток
var (
	TotalRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "Общее количество запросов.",
		},
		[]string{"method", "url", "status"},
	)
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_request_duration_seconds",
			Help:    "Время выполнения запросов.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "url", "status"},
	)
	ExternalRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "external_requests_total",
			Help: "Количество запросов к внешним ресурсам.",
		},
		[]string{"url", "status"},
	)
	ExternalRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "external_request_duration_seconds",
			Help:    "Время выполнения запросов к внешним ресурсам.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"url", "status"},
	)
)

// RecordRequest записывает метрики для запросов к API
func RecordRequest(method, url, status string, duration time.Duration) {
	TotalRequests.WithLabelValues(method, url, status).Inc()
	RequestDuration.WithLabelValues(method, url, status).Observe(duration.Seconds())
}

// RecordExternalRequest записывает метрики для запросов к внешним ресурсам
func RecordExternalRequest(url, status string, duration time.Duration) {
	ExternalRequests.WithLabelValues(url, status).Inc()
	ExternalRequestDuration.WithLabelValues(url, status).Observe(duration.Seconds())
}
