package grpccontroller

import (
	"context"
	"route256/loms/internal/model"
	lomsGrpc "route256/loms/pkg/api/loms/v1"
)

func (o *LomsController) StocksInfo(ctx context.Context, request *lomsGrpc.StocksInfoRequest) (*lomsGrpc.StocksInfoResponse, error) {
	availableCount, err := o.stockService.GetBySKUAvailableCount(ctx, model.SKUType(request.Sku))
	if err != nil {
		return nil, mapErrorToGRPC(err)
	}
	return &lomsGrpc.StocksInfoResponse{Count: availableCount}, nil
}
