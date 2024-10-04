package productservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	client  http.RoundTripper
	token   string
	baseUrl string
	path    string
}

func NewProductService(client http.RoundTripper, token string, baseUrl string, path string) *ProductService {
	return &ProductService{client: client, token: token, baseUrl: baseUrl, path: path}
}

func (p *ProductService) GetProductInfo(ctx context.Context, sku model.SKU) (*model.Product, error) {

	bodyRq := productGetProductRequest{
		Token: p.token,
		Sku:   strconv.FormatInt(int64(sku), 10),
	}

	bodyRaw, err := json.Marshal(bodyRq)
	if err != nil {
		log.Printf("[productservice] Error marshalling request body: %v", err)
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseUrl+p.path, bytes.NewReader(bodyRaw))
	if err != nil {
		log.Printf("[productservice] Error creating request: %v", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.RoundTrip(req)
	if err != nil {
		log.Printf("[productservice] Error sending request: %v", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			log.Printf("[productservice] Error reading error response body for SKU %d: %v", sku, readErr)
			return nil, fmt.Errorf("failed to read error response body: %w", readErr)
		}

		log.Printf("[productservice] Received error response with status %d for SKU %d: %s", resp.StatusCode, sku, string(body))

		switch {
		case resp.StatusCode == http.StatusNotFound:
			return nil, fmt.Errorf("product not found for SKU %d", sku)
		case resp.StatusCode >= 500:
			return nil, fmt.Errorf("server error with status %d for SKU %d: %s", resp.StatusCode, sku, string(body))
		default:
			return nil, fmt.Errorf("unexpected error response with status %d for SKU %d: %s", resp.StatusCode, sku, string(body))
		}
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil || len(content) == 0 {
		log.Printf("[productservice] Error reading response body for SKU %d: %v", sku, err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	err = ctx.Err()
	if err != nil {
		return nil, err
	}
	var productDto productGetProductResponse
	if err := json.Unmarshal(content, &productDto); err != nil {
		log.Printf("[productservice] Error unmarshalling response for SKU %d: %v", sku, err)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	err = ctx.Err()
	if err != nil {
		return nil, err
	}
	if validationErr := validationservice.Validate(productDto); validationErr != nil {
		log.Printf("[productservice] Product validation failed for SKU %d: %v", sku, validationErr)
		return nil, fmt.Errorf("validation failed for SKU %d: %w", sku, validationErr)
	}

	log.Printf("[productservice] Successfully retrieved product info for SKU %d", sku)
	return &model.Product{
		Name:  productDto.Name,
		Price: productDto.Price,
	}, nil
}
