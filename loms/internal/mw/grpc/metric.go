package grpc

import (
	"context"
	"route256/loms/internal/metrics"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func MetricUnaryMiddleware(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	// Засекаем время начала запроса
	start := time.Now()

	// Выполняем обработчик
	resp, err = handler(ctx, req)

	// Определяем статус кода из ошибки, если она есть
	grpcStatus := status.Code(err).String()

	// Длительность запроса
	duration := time.Since(start)

	// Запись метрик
	metrics.RecordRequest("grpc", info.FullMethod, grpcStatus, duration)

	return resp, err
}
