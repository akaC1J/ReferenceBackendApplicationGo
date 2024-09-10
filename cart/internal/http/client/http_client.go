package client

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// RetryRoundTripper реализует интерфейс http.RoundTripper и добавляет логику ретраев.
type RetryRoundTripper struct {
	transport  http.RoundTripper
	maxRetries int
	retryDelay time.Duration
}

func NewHttpClientWithRetry(roundTripper http.RoundTripper, maxRetries int, retryDelayMs int64) *http.Client {
	return &http.Client{
		Transport: &RetryRoundTripper{
			transport:  roundTripper,
			maxRetries: maxRetries,
			retryDelay: time.Duration(retryDelayMs) * time.Millisecond,
		},
	}
}

// NewHttpClientWithRetryWithDefaultTransport Создание HTTP клиента с кастомным RoundTripper (с поддержкой ретраев).
func NewHttpClientWithRetryWithDefaultTransport(maxRetries int, retryDelay time.Duration) *http.Client {
	return &http.Client{
		Transport: &RetryRoundTripper{
			transport:  http.DefaultTransport,
			maxRetries: maxRetries,
			retryDelay: retryDelay,
		},
	}
}

func (rt *RetryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	// Сохраняем тело запроса (если оно не пустое)
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("reading body: %v", err)
		}
	}
	for attempt := 1; attempt <= rt.maxRetries; attempt++ {
		if bodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}
		resp, err = rt.transport.RoundTrip(req)
		//знаю что в спецификации указаны только 429 и 420, но чтобы проверить дейсв
		const httpStatusClientCustomError = 420
		if err == nil && (resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusTooManyRequests &&
			resp.StatusCode != httpStatusClientCustomError) {
			return resp, nil
		}

		if attempt == rt.maxRetries {
			log.Printf("[httpclient] reached the maximum number (%d) of retrays on request: %s\n", rt.maxRetries, req.URL.Host+req.URL.Path)
			break
		}

		log.Println("[httpclient] error. the request will be repeated, url: " + req.URL.Host + req.URL.Path)
		time.Sleep(rt.retryDelay)
	}

	return resp, err
}
