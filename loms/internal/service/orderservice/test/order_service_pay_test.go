package test

import (
	"context"
	"errors"
	"route256/loms/internal/model"
	"route256/loms/internal/service/orderservice"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestService_OrderPay_Success(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()

	repoMock := NewRepositoryMock(mc)
	stockServiceMock := NewStockServiceMock(mc)

	order := &model.Order{
		Items: []*model.Item{{SKU: 1, Count: 10}},
	}

	_ = order.SetState(model.AWAITING_PAYMENT)

	repoMock.GetByIdMock.Expect(ctx, order.ID).Return(order, nil)
	stockServiceMock.ReserveRemoveMock.Expect(ctx, order.Items).Return(nil)
	orderForUpdate := &model.Order{
		ID:    order.ID,
		Items: order.Items,
	}
	_ = orderForUpdate.SetState(model.PAYED)
	repoMock.UpdateOrderMock.Expect(ctx, orderForUpdate).Return(nil)

	service := orderservice.NewService(repoMock, stockServiceMock)

	err := service.OrderPay(ctx, order.ID)
	assert.NoError(t, err)
}

func TestService_OrderPay_GetByIdError(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()
	orderID := int64(1)

	repoMock := NewRepositoryMock(mc)
	stockServiceMock := NewStockServiceMock(mc)

	repoMock.GetByIdMock.Expect(ctx, orderID).Return(nil, errors.New("database error"))

	service := orderservice.NewService(repoMock, stockServiceMock)

	err := service.OrderPay(ctx, orderID)
	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
}

func TestService_OrderPay_ReserveRemoveError(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()

	repoMock := NewRepositoryMock(mc)
	stockServiceMock := NewStockServiceMock(mc)
	order := &model.Order{
		Items: []*model.Item{{SKU: 1, Count: 10}},
	}
	_ = order.SetState(model.AWAITING_PAYMENT)

	repoMock.GetByIdMock.Expect(ctx, order.ID).Return(order, nil)
	stockServiceMock.ReserveRemoveMock.Expect(ctx, order.Items).Return(errors.New("reserve remove error"))

	service := orderservice.NewService(repoMock, stockServiceMock)

	err := service.OrderPay(ctx, order.ID)
	assert.Error(t, err)
	assert.EqualError(t, err, "reserve remove error")
}

func TestService_OrderPay_UpdateOrderError(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()

	repoMock := NewRepositoryMock(mc)
	stockServiceMock := NewStockServiceMock(mc)
	order := &model.Order{
		Items: []*model.Item{{SKU: 1, Count: 10}},
	}
	_ = order.SetState(model.AWAITING_PAYMENT)
	repoMock.GetByIdMock.Expect(ctx, order.ID).Return(order, nil)
	stockServiceMock.ReserveRemoveMock.Expect(ctx, order.Items).Return(nil)
	orderForUpdate := &model.Order{
		ID:    order.ID,
		Items: order.Items,
	}
	_ = orderForUpdate.SetState(model.PAYED)

	repoMock.UpdateOrderMock.Expect(ctx, orderForUpdate).Return(errors.New("update error"))

	service := orderservice.NewService(repoMock, stockServiceMock)

	err := service.OrderPay(ctx, order.ID)
	assert.Error(t, err)
	assert.EqualError(t, err, "update error")
}
