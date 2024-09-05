package client

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

// RetryRoundTripper реализует интерфейс http.RoundTripper и добавляет логику ретраев.
type RetryRoundTripper struct {
	Transport  http.RoundTripper
	MaxRetries int
	RetryDelay time.Duration
}

func (rt *RetryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	// Сохраняем тело запроса (если оно не пустое)
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
	}
	for attempt := 0; attempt <= rt.MaxRetries; attempt++ {
		if bodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}
		resp, err = rt.Transport.RoundTrip(req)
		//знаю что в спецификации указаны только 429 и 420, но чтобы проверить дейсв
		if err == nil && (resp.StatusCode != 404 && resp.StatusCode != 429 && resp.StatusCode != 420) {
			return resp, nil
		}

		if attempt == rt.MaxRetries {
			log.Printf("[httpclient] reached the maximum number (%d) of retrays on request: %s\n", rt.MaxRetries, req.URL.Host+req.URL.Path)
			break
		}

		log.Println("[httpclient] error. the request will be repeated, url: " + req.URL.Host + req.URL.Path)
		time.Sleep(rt.RetryDelay)
	}

	return resp, err
}

// NewHttpClientWithRetry Создание HTTP клиента с кастомным RoundTripper (с поддержкой ретраев).
func NewHttpClientWithRetry(maxRetries int, retryDelayMs int64) *http.Client {
	return &http.Client{
		Transport: &RetryRoundTripper{
			Transport:  http.DefaultTransport,
			MaxRetries: maxRetries,
			RetryDelay: time.Duration(retryDelayMs) * time.Millisecond,
		},
	}
}
