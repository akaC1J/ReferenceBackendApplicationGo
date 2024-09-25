package grpccontroller

import (
	"context"
	lomsGrpc "route256/loms/internal/generated/api/loms/v1"
	"route256/loms/internal/model"
)

func (o *LomsController) StocksInfo(ctx context.Context, request *lomsGrpc.StocksInfoRequest) (*lomsGrpc.StocksInfoResponse, error) {
	availableCount, err := o.stockService.GetBySKUAvailableCount(ctx, model.SKUType(request.Sku))
	if err != nil {
		return nil, mapErrorToGRPC(err)
	}
	return &lomsGrpc.StocksInfoResponse{Count: availableCount}, nil
}
