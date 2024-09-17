package server

import (
	"route256/cart/internal/pkg/model"
	"route256/cart/internal/pkg/service/cartservice"
)

type PostItemRequest struct {
	Count uint16 `json:"count"`
}

type GetCartContentResponse struct {
	*cartservice.CartContent
}

type PostCheckoutRq struct {
	UserId model.UserId `json:"user"`
}

type PostCheckoutRs struct {
	OrderId int64 `json:"orderID"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
