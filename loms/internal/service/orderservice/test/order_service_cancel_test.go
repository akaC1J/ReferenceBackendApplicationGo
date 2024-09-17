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

func TestService_OrderCancel_Success(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()

	repoMock := NewRepositoryMock(mc)
	stockServiceMock := NewStockServiceMock(mc)
	order := &model.Order{
		ID:    1,
		State: model.AWAITING_PAYMENT,
		Items: []*model.Item{{SKU: 1, Count: 10}},
	}

	repoMock.GetByIdMock.Expect(ctx, order.ID).Return(order, nil)
	stockServiceMock.ReserveCancelMock.Expect(ctx, order.Items).Return(nil)

	orderForUpdate := &model.Order{
		ID:    order.ID,
		State: model.CANCELLED,
		Items: order.Items,
	}
	repoMock.UpdateOrderMock.Expect(ctx, orderForUpdate).Return(nil)

	service := orderservice.NewService(repoMock, stockServiceMock)

	err := service.OrderCancel(ctx, order.ID)
	assert.NoError(t, err)
}

func TestService_OrderCancel_GetByIdError(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()
	orderID := int64(1)

	repoMock := NewRepositoryMock(mc)
	stockServiceMock := NewStockServiceMock(mc)

	repoMock.GetByIdMock.Expect(ctx, orderID).Return(nil, errors.New("database error"))

	service := orderservice.NewService(repoMock, stockServiceMock)

	err := service.OrderCancel(ctx, orderID)
	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
}

func TestService_OrderCancel_ReserveCancelError(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()

	repoMock := NewRepositoryMock(mc)
	stockServiceMock := NewStockServiceMock(mc)
	order := &model.Order{
		ID:    1,
		State: model.AWAITING_PAYMENT,
		Items: []*model.Item{{SKU: 1, Count: 10}},
	}
	repoMock.GetByIdMock.Expect(ctx, order.ID).Return(order, nil)
	stockServiceMock.ReserveCancelMock.Expect(ctx, order.Items).Return(errors.New("reserve cancel error"))

	service := orderservice.NewService(repoMock, stockServiceMock)

	err := service.OrderCancel(ctx, order.ID)
	assert.Error(t, err)
	assert.EqualError(t, err, "reserve cancel error")
}

func TestService_OrderCancel_UpdateOrderError(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()

	repoMock := NewRepositoryMock(mc)
	stockServiceMock := NewStockServiceMock(mc)
	order := &model.Order{
		ID:    1,
		State: model.AWAITING_PAYMENT,
		Items: []*model.Item{{SKU: 1, Count: 10}},
	}
	repoMock.GetByIdMock.Expect(ctx, order.ID).Return(order, nil)
	stockServiceMock.ReserveCancelMock.Expect(ctx, order.Items).Return(nil)
	orderForUpdate := &model.Order{
		ID:    order.ID,
		State: model.CANCELLED,
		Items: order.Items,
	}
	repoMock.UpdateOrderMock.Expect(ctx, orderForUpdate).Return(errors.New("update error"))

	service := orderservice.NewService(repoMock, stockServiceMock)

	err := service.OrderCancel(ctx, order.ID)
	assert.Error(t, err)
	assert.EqualError(t, err, "update error")
}
