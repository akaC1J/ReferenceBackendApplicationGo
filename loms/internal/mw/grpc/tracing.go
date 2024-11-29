package grpc

import (
	"context"
	"go.opentelemetry.io/otel/codes"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// TracingUnaryMiddleware извлекает трейсинг-контекст из gRPC-запроса и создаёт новый span
func TracingUnaryMiddleware(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	// Инициализируем Tracer
	tracer := otel.Tracer("grpc-server")

	// Извлекаем метаданные из контекста gRPC
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	// Используем метаданные в качестве Carrier для контекста трейсинга
	propagator := otel.GetTextMapPropagator()
	ctx = propagator.Extract(ctx, metadataCarrier(md))

	// Создаём новый span, продолжая трейсинг, если контекст уже содержит trace_id
	ctx, span := tracer.Start(ctx, info.FullMethod, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	// Передаём обновлённый контекст с span в основной обработчик запроса
	resp, err = handler(ctx, req)

	// Если возникла ошибка, устанавливаем статус span
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "OK")
	}

	return resp, err
}

// metadataCarrier реализует интерфейс TextMapCarrier для gRPC метаданных
type metadataCarrier metadata.MD

func (mc metadataCarrier) Keys() []string {
	keys := make([]string, 0, len(mc))
	for k := range mc {
		keys = append(keys, k)
	}
	return keys
}

func (mc metadataCarrier) Get(key string) string {
	values := metadata.MD(mc).Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (mc metadataCarrier) Set(key string, value string) {
	metadata.MD(mc).Set(key, value)
}
