package tracing

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func InitTracerProvider(serviceName string) (*trace.TracerProvider, error) {
	// Настраиваем экспортёр для Jaeger
	exp, err := otlptracehttp.New(context.Background(), otlptracehttp.WithInsecure(), otlptracehttp.WithEndpoint("localhost:4318"))
	if err != nil {
		return nil, err
	}

	// Создаем ресурс для указания информации о сервисе
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	// Создаём TracerProvider с пропагацией контекста
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(res),
	)

	// Устанавливаем TracerProvider как глобальный и добавляем пропагатор
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{}) // используем W3C Trace Context

	return tp, nil
}
