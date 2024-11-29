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
	DBRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_requests_total",
			Help: "Количество запросов в базу данных по категории (select, update, delete).",
		},
		[]string{"category"},
	)
	DBRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_request_duration_seconds",
			Help:    "Время выполнения запросов в базу данных по категории.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"category", "status"},
	)
)

// RecordRequest записывает метрики для запросов к API
func RecordRequest(method, url, status string, duration time.Duration) {
	TotalRequests.WithLabelValues(method, url, status).Inc()
	RequestDuration.WithLabelValues(method, url, status).Observe(duration.Seconds())
}

// RecordDBRequest записывает метрики для запросов к базе данных
func RecordDBRequest(category, status string, duration time.Duration) {
	DBRequests.WithLabelValues(category).Inc()
	DBRequestDuration.WithLabelValues(category, status).Observe(duration.Seconds())
}
