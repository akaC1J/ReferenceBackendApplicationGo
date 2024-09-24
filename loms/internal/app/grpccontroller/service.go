package grpccontroller

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	appErr "route256/loms/internal/errors"
	"route256/loms/internal/model"
	lomsGrpc "route256/loms/pkg/api/loms/v1"
)

var _ lomsGrpc.LomsServer = (*LomsController)(nil)

type OrderService interface {
	Create(ctx context.Context, order *model.Order) (orderID int64, err error)
	GetById(ctx context.Context, orderID int64) (*model.Order, error)
	OrderPay(ctx context.Context, orderID int64) error
	OrderCancel(ctx context.Context, orderID int64) error
}

type StockService interface {
	GetBySKUAvailableCount(ctx context.Context, sku model.SKUType) (uint64, error)
}
type LomsController struct {
	orderService OrderService
	stockService StockService
	lomsGrpc.UnimplementedLomsServer
}

func NewLomsController(orderService OrderService, stockService StockService) *LomsController {
	return &LomsController{orderService: orderService, stockService: stockService}
}

func mapErrorToGRPC(err error) error {
	if errors.Is(err, appErr.ErrStockInsufficient) {
		return status.Error(codes.FailedPrecondition, err.Error())
	}
	if errors.Is(err, appErr.ErrNotFound) || errors.Is(err, appErr.ErrOrderState) {
		return status.Error(codes.NotFound, err.Error())
	}
	if errors.Is(err, appErr.ErrNegativeReserved) || errors.Is(err, appErr.ErrNegativeAvailable) {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return status.Error(codes.Internal, "Internal server error")
}
