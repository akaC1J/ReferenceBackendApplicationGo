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

func TestService_Create_Success(t *testing.T) {
	mc := minimock.NewController(t)
	ctx := context.Background()

	repoMock := NewRepositoryMock(mc)
	stockServiceMock := NewStockServiceMock(mc)
	order := &model.Order{
		Items: []*model.Item{{SKU: 1, Count: 10}},
	}

	savedOrder := &model.Order{
		Items: order.Items,
		ID:    1,
	}
	repoMock.SaveOrderMock.Expect(ctx, order).Return(savedOrder, nil)
	stockServiceMock.ReserveMock.Expect(ctx, order.Items).Return(nil)
	_ = savedOrder.SetState(model.AWAITING_PAYMENT)
	repoMock.UpdateOrderMock.Expect(ctx, savedOrder).Return(nil)

	service := orderservice.NewService(repoMock, stockServiceMock)

	orderID, err := service.Create(ctx, order)
	assert.NoError(t, err)
	assert.Equal(t, savedOrder.ID, orderID)
}

func TestService_Create_SaveOrderError(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()

	repoMock := NewRepositoryMock(mc)
	stockServiceMock := NewStockServiceMock(mc)

	order := &model.Order{
		Items: []*model.Item{{SKU: 1, Count: 10}},
	}

	repoMock.SaveOrderMock.Expect(ctx, order).Return(nil, errors.New("save error"))

	service := orderservice.NewService(repoMock, stockServiceMock)

	orderID, err := service.Create(ctx, order)
	assert.Error(t, err)
	assert.Equal(t, int64(0), orderID)
}

func TestService_Create_ReserveError(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()

	order := &model.Order{
		Items: []*model.Item{{SKU: 1, Count: 10}},
	}
	savedOrder := &model.Order{
		Items: order.Items,
		ID:    1,
	}

	repoMock := NewRepositoryMock(mc)
	stockServiceMock := NewStockServiceMock(mc)

	repoMock.SaveOrderMock.Expect(ctx, order).Return(savedOrder, nil)
	stockServiceMock.ReserveMock.Expect(ctx, order.Items).Return(errors.New("reserve error"))
	_ = savedOrder.SetState(model.FAILED)

	repoMock.UpdateOrderMock.Expect(ctx, savedOrder).Return(nil)

	service := orderservice.NewService(repoMock, stockServiceMock)

	orderID, err := service.Create(ctx, order)
	assert.Error(t, err)
	assert.Equal(t, int64(0), orderID)
}

func TestService_Create_UpdateOrderError(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()

	order := &model.Order{
		Items: []*model.Item{{SKU: 1, Count: 10}},
	}

	savedOrder := &model.Order{
		Items: order.Items,
		ID:    1,
	}

	repoMock := NewRepositoryMock(mc)
	stockServiceMock := NewStockServiceMock(mc)

	repoMock.SaveOrderMock.Expect(ctx, order).Return(savedOrder, nil)
	stockServiceMock.ReserveMock.Expect(ctx, order.Items).Return(nil)
	_ = savedOrder.SetState(model.AWAITING_PAYMENT)

	repoMock.UpdateOrderMock.Expect(ctx, savedOrder).Return(errors.New("update error"))

	service := orderservice.NewService(repoMock, stockServiceMock)

	orderID, err := service.Create(ctx, order)
	assert.Error(t, err)
	assert.Equal(t, int64(0), orderID)
}
