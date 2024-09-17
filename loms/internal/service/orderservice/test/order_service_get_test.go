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

func TestService_GetById_Success(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()

	order := &model.Order{
		Items: []*model.Item{{SKU: 1, Count: 10}},
	}

	repoMock := NewRepositoryMock(mc)
	repoMock.GetByIdMock.Expect(ctx, order.ID).Return(order, nil)

	service := orderservice.NewService(repoMock, NewStockServiceMock(mc))

	order, err := service.GetById(ctx, order.ID)
	assert.NoError(t, err)
	assert.Equal(t, order, order)
}

func TestService_GetById_Error(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()
	orderID := int64(1)

	repoMock := NewRepositoryMock(mc)
	repoMock.GetByIdMock.Expect(ctx, orderID).Return(nil, errors.New("database error"))

	service := orderservice.NewService(repoMock, NewStockServiceMock(mc))

	order, err := service.GetById(ctx, orderID)
	assert.Error(t, err)
	assert.Nil(t, order)
	assert.EqualError(t, err, "database error")
}
