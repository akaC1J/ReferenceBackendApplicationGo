package server

import "route256/cart/internal/pkg/service/cartservice"

type PostItemRequest struct {
	Count uint16 `json:"count"`
}

type GetCartContentResponse struct {
	*cartservice.CartContent
}

type ErrorResponse struct {
	Message string `json:"message"`
}
