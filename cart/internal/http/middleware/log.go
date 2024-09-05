package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type LogMux struct {
	h http.Handler
}

func NewLogMux(h http.Handler) http.Handler {
	return &LogMux{h: h}
}

type responseCapture struct {
	http.ResponseWriter
	body *bytes.Buffer
	code int
}

func (rc *responseCapture) WriteHeader(statusCode int) {
	rc.code = statusCode
	rc.ResponseWriter.WriteHeader(statusCode)
}

func (rc *responseCapture) Write(b []byte) (int, error) {
	rc.body.Write(b)
	return rc.ResponseWriter.Write(b)
}

func (m *LogMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	clientIP := r.RemoteAddr
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		clientIP = ip
	}

	rqMessageLog := fmt.Sprintf("Request: \nMethod: %s\nPath: %s\nClientIp: %s\n", r.Method, r.URL.Path, clientIP)

	if r.Body != nil {
		bodyBytes, err := io.ReadAll(r.Body)
		if err == nil {
			body := strings.TrimSpace(string(bodyBytes))
			if len(body) > 0 {
				rqMessageLog += fmt.Sprintf("Body: %s\n", body)
			}
		}
		r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
	}

	rc := &responseCapture{
		ResponseWriter: w,
		body:           bytes.NewBuffer(nil),
	}

	m.h.ServeHTTP(rc, r)

	responseLog := fmt.Sprintf("Response: \nStatusCode: %d\nBody: %s", rc.code, rc.body.String())

	duration := time.Since(startTime)

	fullLog := fmt.Sprintf("%s%s\nDuration: %v", rqMessageLog, responseLog, duration)
	log.Println("[middleware] " + fullLog)
}
