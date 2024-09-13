package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"route256/cart/internal/app"
	"route256/cart/internal/pkg/service/cartservice"
	"strconv"
	"testing"
	"time"

	"route256/cart/internal/pkg/model"
)

func startTestServer(stopChan chan struct{}) {
	os.Setenv("ENV", "TEST")
	config, err := app.LoadConfig()
	if err != nil {
		log.Fatalf("[main] Failed to load configuration: %v", err)
	}

	application, err := app.New(config)
	if err != nil {
		log.Fatalf("[main] Failed to initialize application: %v", err)
	}

	server := &http.Server{
		Addr:    config.Host_Port,
		Handler: application,
	}

	go func() {
		<-stopChan
		log.Println("[main] Stopping server...")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatalf("[main] Server shutdown failed: %v", err)
		}
	}()

	log.Printf("[main] Starting server on %s\n", config.Host_Port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("[main] Failed to start server: %v", err)
	}
}
func TestGetCartContent_E2E(t *testing.T) {
	stopChan := make(chan struct{})
	go startTestServer(stopChan)
	// Ждём, пока сервер запустится
	time.Sleep(2 * time.Second) // Возможно, потребуется увеличить время ожидания

	// Подготовка тестовых данных
	userID := int64(123)
	sku := int64(1076963)
	item := model.CartItem{
		SKU:   model.SKU(sku),
		Count: 2,
	}

	// Отправляем POST запрос для добавления товара в корзину
	addItemURL := "http://localhost:8082/user/" + strconv.FormatInt(userID, 10) + "/cart/" + strconv.FormatInt(sku, 10)
	addItemBody, _ := json.Marshal(map[string]interface{}{
		"count": item.Count,
	})
	req, err := http.NewRequest("POST", addItemURL, bytes.NewBuffer(addItemBody))
	if err != nil {
		t.Fatalf("Не удалось создать POST запрос: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Не удалось выполнить POST запрос: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Ожидался статус код %d при добавлении товара, получен %d", http.StatusOK, resp.StatusCode)
	}

	// Отправляем GET запрос для получения содержимого корзины
	getCartURL := "http://localhost:8082/user/" + strconv.FormatInt(userID, 10) + "/cart"
	resp, err = http.Get(getCartURL)
	if err != nil {
		t.Fatalf("Не удалось выполнить GET запрос: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем код статуса
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Ожидался статус код %d, получен %d", http.StatusOK, resp.StatusCode)
	}

	// Декодируем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Не удалось прочитать тело ответа: %v", err)
	}

	var response struct {
		Items      []cartservice.EnrichedCartItem `json:"items"`
		TotalPrice uint32                         `json:"total_price"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("Не удалось декодировать тело ответа: %v", err)
	}

	// Проверяем содержимое ответа
	if len(response.Items) != 1 {
		t.Fatalf("Ожидался 1 товар в корзине, получено %d", len(response.Items))
	}

	itemResp := response.Items[0]
	if itemResp.SKU != sku {
		t.Errorf("Ожидался SKU %d, получен %d", sku, itemResp.SKU)
	}
	if itemResp.Count != item.Count {
		t.Errorf("Ожидалось количество %d, получено %d", item.Count, itemResp.Count)
	}

	deleteItemURL := "http://localhost:8082/user/" + strconv.FormatInt(userID, 10) + "/cart/" + strconv.FormatInt(sku, 10)
	req, err = http.NewRequest("DELETE", deleteItemURL, nil)
	if err != nil {
		t.Fatalf("Не удалось создать DELETE запрос: %v", err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Не удалось выполнить DELETE запрос: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Ожидался статус код %d при удалении товара, получен %d", http.StatusNoContent, resp.StatusCode)
	}

	close(stopChan)
}
