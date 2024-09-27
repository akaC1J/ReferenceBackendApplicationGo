package grpccontroller

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	lomsGrpc "route256/loms/internal/generated/api/loms/v1"
	"route256/loms/internal/model"
)

func (o *LomsController) OrderCreate(ctx context.Context, createRq *lomsGrpc.OrderCreateRequest) (*lomsGrpc.OrderCreateResponse, error) {
	orderId, err := o.orderService.Create(ctx, convertCreateRequestToOrder(createRq))
	if err != nil {
		return nil, mapErrorToGRPC(err)
	}
	return &lomsGrpc.OrderCreateResponse{OrderId: orderId}, nil
}

func (o *LomsController) OrderPay(ctx context.Context, request *lomsGrpc.OrderPayRequest) (*emptypb.Empty, error) {
	err := o.orderService.OrderPay(ctx, request.OrderId)
	if err != nil {
		return nil, mapErrorToGRPC(err)
	}
	return &emptypb.Empty{}, nil
}

func (o *LomsController) OrderCancel(ctx context.Context, request *lomsGrpc.OrderCancelRequest) (*emptypb.Empty, error) {
	err := o.orderService.OrderCancel(ctx, request.OrderId)
	if err != nil {
		return nil, mapErrorToGRPC(err)
	}
	return &emptypb.Empty{}, nil
}

func (o *LomsController) OrderInfo(ctx context.Context, request *lomsGrpc.OrderInfoRequest) (*lomsGrpc.OrderInfoResponse, error) {
	order, err := o.orderService.GetById(ctx, request.OrderId)
	if err != nil {
		return nil, mapErrorToGRPC(err)
	}
	return &lomsGrpc.OrderInfoResponse{Order: convertOrderToResponse(order)}, nil
}

func convertOrderToResponse(order *model.Order) *lomsGrpc.Order {
	orderRs := &lomsGrpc.Order{
		User: order.UserId,
	}
	for _, item := range order.Items {
		orderRs.Items = append(orderRs.Items, &lomsGrpc.Item{
			Sku:   uint32(item.SKU),
			Count: item.Count,
		})
	}

	orderRs.Status = string(order.State())

	return orderRs
}

func convertCreateRequestToOrder(createRq *lomsGrpc.OrderCreateRequest) *model.Order {
	order := model.Order{
		UserId: createRq.GetOrder().GetUser(),
	}
	for _, itemRq := range createRq.GetOrder().GetItems() {
		order.Items = append(order.Items, &model.Item{
			SKU:   model.SKUType(itemRq.GetSku()),
			Count: itemRq.GetCount(),
		})
	}
	return &order
}
