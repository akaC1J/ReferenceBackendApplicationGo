//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	app "route256/loms/internal/app/initialization"
	lomsGrpc "route256/loms/pkg/api/loms/v1" // Путь к сгенерированным gRPC клиентам
)

func TestE2E_OrderLifecycle(t *testing.T) {

	//Инициализация и запуск сервера
	config, err := app.LoadConfig("./.env.test")
	assert.NoError(t, err, "Не удалось загрузить конфигурацию")

	application, err := app.New(config)
	assert.NoError(t, err, "Не удалось инициализировать приложение")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GrpcPort))
	assert.NoError(t, err, "Не удалось создать слушатель для gRPC сервера")

	go func() {
		if err := application.GrpcServer.Serve(lis); err != nil {
			log.Fatalf("[e2e_test] gRPC сервер завершился с ошибкой: %v", err)
		}
	}()
	defer application.GrpcServer.GracefulStop()

	// Подождем немного, чтобы сервер успел запуститься
	time.Sleep(2 * time.Second)

	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", config.GrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "Не удалось подключиться к gRPC серверу")
	defer conn.Close()

	client := lomsGrpc.NewLomsClient(conn)

	ctx := context.Background()

	//Выполнение OrderCreate запроса
	orderCreateReq := &lomsGrpc.OrderCreateRequest{
		Order: &lomsGrpc.Order{
			User: 12345,
			Items: []*lomsGrpc.Item{
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

	//Выполнение первого OrderInfo запроса
	orderInfoReq1 := &lomsGrpc.OrderInfoRequest{
		OrderId: orderID,
	}

	orderInfoResp1, err := client.OrderInfo(ctx, orderInfoReq1)
	assert.NoError(t, err, "OrderInfo запрос (1) завершился с ошибкой")
	assert.NotNil(t, orderInfoResp1, "OrderInfo ответ (1) пустой")
	assert.Equal(t, "awaiting_payment", orderInfoResp1.Order.Status, "Статус заказа не соответствует ожидаемому (AWAITING_PAY)")

	// Шаг 7: Выполнение OrderPay запроса
	orderPayReq := &lomsGrpc.OrderPayRequest{
		OrderId: orderID,
	}

	_, err = client.OrderPay(ctx, orderPayReq)
	assert.NoError(t, err, "OrderPay запрос завершился с ошибкой")

	//Выполнение второго OrderInfo запроса
	orderInfoReq2 := &lomsGrpc.OrderInfoRequest{
		OrderId: orderID,
	}

	orderInfoResp2, err := client.OrderInfo(ctx, orderInfoReq2)
	assert.NoError(t, err, "OrderInfo запрос (2) завершился с ошибкой")
	assert.NotNil(t, orderInfoResp2, "OrderInfo ответ (2) пустой")
	assert.Equal(t, "payed", orderInfoResp2.Order.Status, "Статус заказа не соответствует ожидаемому (PAY)")

	//Выполнение StocksInfo запроса
	stocksInfoReq := &lomsGrpc.StocksInfoRequest{
		Sku: 1,
	}

	stocksInfoResp, err := client.StocksInfo(ctx, stocksInfoReq)
	assert.NoError(t, err, "StocksInfo запрос завершился с ошибкой")
	assert.NotNil(t, stocksInfoResp, "StocksInfo ответ пустой")
	expectedCount := uint64(135) // Начальные запасы 140, было заказано 5
	assert.Equal(t, expectedCount, stocksInfoResp.Count, "Количество запасов SKU 1002 не соответствует ожидаемому")
}
