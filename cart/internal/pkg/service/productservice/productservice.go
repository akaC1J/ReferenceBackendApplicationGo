package productservice

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"route256/cart/internal/pkg/model"
	"route256/cart/internal/pkg/service/validationservice"
	"strconv"
)

type ProductService struct {
	// решил здесь принебречь интерфейсом, потому что считаю, что http клиент это что-то конкретное,
	// а не абстракцию
	client  *http.Client
	token   string
	baseUrl string
	path    string
}

func NewProductService(client *http.Client, token string, baseUrl string, path string) *ProductService {
	return &ProductService{client: client, token: token, baseUrl: baseUrl, path: path}
}

func (p *ProductService) GetProductInfo(_ context.Context, sku model.SKU) (*model.Product, error) {

	// Подготовка запроса
	bodyRq := productGetProductRequest{
		Token: p.token,
		Sku:   strconv.FormatInt(int64(sku), 10),
	}
	bodyRaw, err := json.Marshal(bodyRq)
	if err != nil {
		log.Printf("[productservice] Error marshalling request body: %v", err)
		return nil, err
	}

	resp, err := p.client.Post(p.baseUrl+p.path, "application/json", bytes.NewReader(bodyRaw))
	if err != nil {
		log.Printf("[productservice] Error sending request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Проверка на успешность ответа
	if resp.StatusCode >= 400 {
		log.Printf("[productservice] Received error response: %d", resp.StatusCode)
		strCont, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[productservice] Error reading response body: %v", err)
			return nil, err
		}
		log.Printf("[productservice] Error message from server: %s", string(strCont))
		return nil, errors.New(string(strCont))
	}

	cont, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[productservice] Error reading response body: %v", err)
		return nil, err
	}

	productDto := new(productGetProductResponse)
	err = json.Unmarshal(cont, productDto)
	if err != nil {
		log.Printf("[productservice] Error unmarshalling response to Product struct: %v", err)
		return nil, err
	}

	validationErr := validationservice.Validate(productDto)
	if validationErr != nil {
		log.Printf("[productservice] Product validation failed for SKU %d: %v", sku, validationErr)
		return nil, validationErr
	}

	// Успешное получение информации о продукте
	log.Printf("[productservice] Successfully retrieved productDto info for SKU %d", sku)
	log.Printf("[productservice] productDto: %+v", productDto)
	return &model.Product{
		Name:  productDto.Name,
		Price: productDto.Price,
	}, nil
}
