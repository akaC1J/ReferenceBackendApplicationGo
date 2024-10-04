package productservice

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"route256/cart/internal/pkg/model"
	"testing"
)

type ProductServiceSuite struct {
	suite.Suite
	server     *httptest.Server
	productSvc *ProductService
}

const (
	defaultOkStatusPath = "/status_ok"
	notFoundStatusFound = "/status_notfound"
)

// TestMain запускает все тесты из структуры ProductServiceSuite
// оставим этот suite последовательным, потому что многие тесты конкурирует за объект сервиса, и некоторые
// тесты могут быть не стабильными, так как изменяют поведения общего ресурса
func TestProductServiceSuite(t *testing.T) {
	suite.Run(t, new(ProductServiceSuite))
}

// SetupSuite выполняется один раз перед всеми тестами
func (suite *ProductServiceSuite) SetupTest() {
	// Настраиваем поддельный сервер для успешного ответа
	suite.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == defaultOkStatusPath {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"name": "Test Product", "price": 100}`))
		} else if r.URL.Path == notFoundStatusFound {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "product not found"}`))
		}
	}))

	client := suite.server.Client().Transport
	suite.productSvc = NewProductService(client, "test-token", suite.server.URL, defaultOkStatusPath)
}

func (suite *ProductServiceSuite) TearDownTest() {
	suite.server.Close()
}

// TestGetProductInfo_Success проверяет успешное получение информации о продукте
func (suite *ProductServiceSuite) TestGetProductInfo_Success() {
	product, err := suite.productSvc.GetProductInfo(context.Background(), model.SKU(123))

	// Проверка результата
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), product)
	assert.Equal(suite.T(), "Test Product", product.Name)
	assert.Equal(suite.T(), uint32(100), product.Price)
}

// TestGetProductInfo_NotFound проверяет, что продукт не найден
func (suite *ProductServiceSuite) TestGetProductInfo_NotFound() {
	suite.productSvc.path = notFoundStatusFound

	product, err := suite.productSvc.GetProductInfo(context.Background(), model.SKU(999))

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), product)
	assert.Contains(suite.T(), err.Error(), "product not found")
}

// TestGetProductInfo_RequestError проверяет ошибку при создании запроса
func (suite *ProductServiceSuite) TestGetProductInfo_RequestError() {
	// Настраиваем продуктовый сервис с некорректным URL, который вызовет ошибку создания запроса
	productSvc := NewProductService(http.DefaultTransport, "test-token", "http://::invalid-url", "/test-path")

	_, err := productSvc.GetProductInfo(context.Background(), model.SKU(123))

	// Проверяем, что ошибка создания запроса
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to create request")
}

// TestGetProductInfo_SendRequestError проверяет ошибку при отправке запроса
func (suite *ProductServiceSuite) TestGetProductInfo_SendRequestError() {
	// Закрываем сервер, чтобы симулировать ошибку отправки запроса
	suite.server.Close()

	_, err := suite.productSvc.GetProductInfo(context.Background(), model.SKU(123))

	// Проверяем, что ошибка отправки запроса
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to send request")
}

// TestGetProductInfo_ReadResponseError проверяет ошибку при чтении тела ответа
func (suite *ProductServiceSuite) TestGetProductInfo_ReadResponseError() {
	// Настраиваем сервер, который не возвращает тело ответа
	suite.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(nil) // Пустое тело ответа
	})

	_, err := suite.productSvc.GetProductInfo(context.Background(), model.SKU(123))

	// Проверяем, что ошибка чтения тела ответа
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to read response body")
}

// TestGetProductInfo_InternalServerError проверяет ошибку при коде ответа 500
func (suite *ProductServiceSuite) TestGetProductInfo_InternalServerError() {
	// Настраиваем сервер для возврата 500 Internal Server Error
	suite.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "internal server error"}`))
	})

	_, err := suite.productSvc.GetProductInfo(context.Background(), model.SKU(123))

	// Проверяем, что ошибка с кодом 500
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "server error with status 500")
}

// TestGetProductInfo_UnmarshalError проверяет ошибку при распарсивании ответа
func (suite *ProductServiceSuite) TestGetProductInfo_UnmarshalError() {
	// Настраиваем сервер для возврата некорректного JSON
	suite.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name": "Test Product", "price": "invalid_price"}`)) // Некорректный JSON
	})

	product, err := suite.productSvc.GetProductInfo(context.Background(), model.SKU(123))

	// Проверяем, что произошла ошибка распарсивания
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), product)
	assert.Contains(suite.T(), err.Error(), "failed to unmarshal response")
}

// TestGetProductInfo_ValidationError проверяет ошибку при валидации ответа
func (suite *ProductServiceSuite) TestGetProductInfo_ValidationError() {
	suite.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"price": 5}`)) // Некорректные данные: пустое имя
	})

	product, err := suite.productSvc.GetProductInfo(context.Background(), model.SKU(123))

	// Проверяем, что произошла ошибка валидации
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), product)
	assert.Contains(suite.T(), err.Error(), "validation failed for SKU")
}
