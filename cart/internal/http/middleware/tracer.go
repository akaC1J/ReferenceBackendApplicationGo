package middleware

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
)

type TraceMux struct {
	h http.Handler
}

// NewTraceMux создаёт новый TraceMux
func NewTraceMux(h http.Handler) http.Handler {
	return &TraceMux{h: h}
}

// ServeHTTP обрабатывает каждый запрос, создавая или продолжая трейс
func (tm *TraceMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Определяем имя span по URL
	operationName := "http_request:" + cleanRequestPath(r.URL.Path)

	// Получаем глобальный Tracer
	tracer := otel.Tracer("trace-mux")

	// Извлекаем контекст трейсинга из заголовков HTTP-запроса
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

	// Создаём новый span или продолжаем существующий трейс, если контекст уже содержит trace_id
	ctx, span := tracer.Start(ctx, operationName)
	defer span.End()

	// Передаем обновленный контекст в следующий обработчик
	r = r.WithContext(ctx)
	tm.h.ServeHTTP(w, r)
}
