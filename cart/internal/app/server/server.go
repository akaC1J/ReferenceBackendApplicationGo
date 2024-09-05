package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"route256/cart/internal/pkg/model"
	"route256/cart/internal/pkg/service/cartservice"
)

type CartInterface interface {
	AddCartItem(ctx context.Context, cartItem model.CartItem) (*model.CartItem, error)
	DeleteCartItem(ctx context.Context, userId model.UserId, sku model.SKU) error
	CleanUpCart(ctx context.Context, userId model.UserId) error
	GetCartItem(ctx context.Context, userId model.UserId) (*cartservice.CartContent, error)
}

type Server struct {
	reviewService CartInterface
}

func New(reviewService CartInterface) *Server {
	return &Server{reviewService: reviewService}
}

func respondWithError(w http.ResponseWriter, statusCode int, message, methodUrl string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	type ErrorResponse struct {
		Message string `json:"message"`
	}
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
