//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"log"
	"net"
	"route256/loms/internal/generated/api/loms/v1"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	app "route256/loms/internal/app/initialization"
)

func TestE2E_OrderCancellationLifecycle(t *testing.T) {
	// Шаг 1: Загрузка конфигурации
	config, err := app.LoadConfig("./.env.test")
	assert.NoError(t, err, "Не удалось загрузить конфигурацию")

	// Шаг 2: Инициализация и запуск приложения
	application, err := app.MustNew(config)
	assert.NoError(t, err, "Не удалось инициализировать приложение")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GgrpcHostPort))
	assert.NoError(t, err, "Не удалось создать слушатель для gRPC сервера")

	go func() {
		if err := application.GrpcServer.Serve(lis); err != nil {
			log.Fatalf("[e2e_test] gRPC сервер завершился с ошибкой: %v", err)
		}
	}()
	defer application.GrpcServer.GracefulStop()

	// Подождём немного, чтобы сервер успел запуститься
	time.Sleep(2 * time.Second)

	// Шаг 3: Создание gRPC клиента
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", config.GgrpcHostPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "Не удалось подключиться к gRPC серверу")
	defer conn.Close()

	client := loms.NewLomsClient(conn)

	ctx := context.Background()

	// Шаг 4: Выполнение OrderCreate запроса
	orderCreateReq := &loms.OrderCreateRequest{
		Order: &loms.Order{
			User: 12345,
			Items: []*loms.Item{
				{Sku: 1, Count: 5},
				{Sku: 2, Count: 10},
			},
		},
	}

	orderCreateResp, err := client.OrderCreate(ctx, orderCreateReq)
	assert.NoError(t, err, "OrderCreate запрос завершился с ошибкой")
	assert.NotNil(t, orderCreateResp, "OrderCreate ответ пустой")
	assert.Greater(t, orderCreateResp.OrderId, int64(0), "Получен некорректный orderId")

	orderID := orderCreateResp.OrderId

	// Шаг 5: Выполнение первого OrderInfo запроса
	orderInfoReq1 := &loms.OrderInfoRequest{
		OrderId: orderID,
	}

	orderInfoResp1, err := client.OrderInfo(ctx, orderInfoReq1)
	assert.NoError(t, err, "OrderInfo запрос (1) завершился с ошибкой")
	assert.NotNil(t, orderInfoResp1, "OrderInfo ответ (1) пустой")
	assert.Equal(t, "awaiting_payment", orderInfoResp1.Order.Status, "Статус заказа не соответствует ожидаемому (AWAITING_PAY)")

	// Шаг 6: Выполнение OrderCancel запроса
	orderCancelReq := &loms.OrderCancelRequest{
		OrderId: orderID,
	}

	_, err = client.OrderCancel(ctx, orderCancelReq)
	assert.NoError(t, err, "OrderCancel запрос завершился с ошибкой")
	// Предполагаем, что OrderCancelResponse содержит подтверждение отмены или просто успешный ответ

	// Шаг 7: Выполнение второго OrderInfo запроса
	orderInfoReq2 := &loms.OrderInfoRequest{
		OrderId: orderID,
	}

	orderInfoResp2, err := client.OrderInfo(ctx, orderInfoReq2)
	assert.NoError(t, err, "OrderInfo запрос (2) завершился с ошибкой")
	assert.NotNil(t, orderInfoResp2, "OrderInfo ответ (2) пустой")
	assert.Equal(t, "cancelled", orderInfoResp2.Order.Status, "Статус заказа не соответствует ожидаемому (CANCELLED)")

	// Шаг 8: Выполнение StocksInfo запроса для SKU 1003
	stocksInfoReq := &loms.StocksInfoRequest{
		Sku: 1,
	}

	stocksInfoResp, err := client.StocksInfo(ctx, stocksInfoReq)
	assert.NoError(t, err, "StocksInfo запрос завершился с ошибкой")
	assert.NotNil(t, stocksInfoResp, "StocksInfo ответ пустой")

	// Предполагаем, что начальные запасы для SKU 1003 были 140
	initialCount := uint64(140)
	expectedCount := initialCount // Поскольку заказ был отменен, запасы должны остаться неизменными

	assert.Equal(t, expectedCount, stocksInfoResp.Count, "Количество запасов SKU 1003 не соответствует ожидаемому")
}
