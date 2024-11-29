package middleware

import (
	"net/http"
	"regexp"
	"route256/cart/internal/metrics"
	"time"
)

type MetricMux struct {
	nextHanler http.Handler
}

func NewMetricMux(nextHandler http.Handler) http.Handler {
	return &MetricMux{nextHanler: nextHandler}
}

type statusRecorder struct {
	http.ResponseWriter
	status string
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.status = http.StatusText(statusCode)
	r.ResponseWriter.WriteHeader(statusCode)
}

func (m *MetricMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	statusRecorderWriter := &statusRecorder{ResponseWriter: w, status: http.StatusText(http.StatusOK)}

	m.nextHanler.ServeHTTP(statusRecorderWriter, r)

	duration := time.Since(start)

	metrics.RecordRequest(r.Method, cleanRequestPath(r.URL.Path), statusRecorderWriter.status, duration)
}

var reUserID = regexp.MustCompile(`/user/[^/]+`)
var reSkuID = regexp.MustCompile(`/cart/[^/]+`)

func cleanRequestPath(path string) string {

	// Заменяем переменные части на шаблоны
	path = reUserID.ReplaceAllString(path, "/user/{user_id}")
	path = reSkuID.ReplaceAllString(path, "/cart/{sku_id}")

	return path
}
