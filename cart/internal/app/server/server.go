package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"route256/cart/internal/pkg/model"
	"route256/cart/internal/pkg/service/cartservice"
	"strconv"
)

type CartInterface interface {
	AddCartItem(ctx context.Context, cartItem model.CartItem) (*model.CartItem, error)
	DeleteCartItem(ctx context.Context, userId model.UserId, sku model.SKU) error
	CleanUpCart(ctx context.Context, userId model.UserId) error
	GetCartItem(ctx context.Context, userId model.UserId) (*cartservice.CartContent, error)
	Checkout(ctx context.Context, userId model.UserId) (orderId int64, err error)
}

type Server struct {
	cartInterface CartInterface
}

func New(reviewService CartInterface) *Server {
	return &Server{cartInterface: reviewService}
}

func respondWithError(w http.ResponseWriter, statusCode int, message, methodUrl string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	// Создаем структуру с ошибкой
	errorResponse := ErrorResponse{
		Message: message,
	}

	// Сериализуем ее в JSON
	err := json.NewEncoder(w).Encode(errorResponse)
	if err != nil {
		log.Printf("Response %s writing failed: %s", methodUrl, err.Error())
	}
}

func getParamFromReq(r *http.Request, paramName string) (int64, error) {
	rawValue := r.PathValue(paramName)
	if rawValue == "" {
		return 0, fmt.Errorf("missing or empty parameter: %s", paramName)
	}

	parsedValue, err := strconv.ParseInt(rawValue, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s format", paramName)
	}

	return parsedValue, nil
}
