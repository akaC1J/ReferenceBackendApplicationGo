package test

import (
	"context"
	"errors"
	appErrors "route256/loms/internal/errors"
	"route256/loms/internal/model"
	"route256/loms/internal/service/stockservice"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestService_ReserveCancel_Success(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()

	items := []*model.Item{
		{SKU: 1, Count: 5},
		{SKU: 2, Count: 3},
	}

	repoMock := NewRepositoryMock(mc)

	repoMock.GetStockMock.
		When(ctx, model.SKUType(1)).
		Then(&model.Stock{SKU: 1, TotalCount: 20, ReservedCount: 10}, nil)

	repoMock.GetStockMock.
		When(ctx, model.SKUType(2)).
		Then(&model.Stock{SKU: 2, TotalCount: 15, ReservedCount: 5}, nil)

	expectedUpdateStocks := map[model.SKUType]*model.Stock{
		1: {SKU: 1, TotalCount: 20, ReservedCount: 5}, // 10 - 5
		2: {SKU: 2, TotalCount: 15, ReservedCount: 2}, // 5 - 3
	}
	repoMock.UpdateStockMock.Set(func(ctx context.Context, stocks map[model.SKUType]*model.Stock) error {
		assert.True(t, compareStocks(expectedUpdateStocks, stocks))
		return nil
	})

	service := stockservice.NewService(repoMock)

	err := service.ReserveCancel(ctx, items)
	assert.NoError(t, err)
}

func TestService_ReserveCancel_GetStockError(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()

	items := []*model.Item{
		{SKU: 1, Count: 5},
	}

	repoMock := NewRepositoryMock(mc)

	repoMock.GetStockMock.
		When(ctx, model.SKUType(1)).
		Then(&model.Stock{}, errors.New("database error"))

	service := stockservice.NewService(repoMock)

	err := service.ReserveCancel(ctx, items)
	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
}

func TestService_ReserveCancel_NegativeReservedCount(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()

	items := []*model.Item{
		{SKU: 1, Count: 15},
	}

	repoMock := NewRepositoryMock(mc)

	repoMock.GetStockMock.
		When(ctx, model.SKUType(1)).
		Then(&model.Stock{SKU: 1, TotalCount: 20, ReservedCount: 10}, nil)

	service := stockservice.NewService(repoMock)

	err := service.ReserveCancel(ctx, items)
	assert.Error(t, err)
	assert.IsType(t, appErrors.ErrStockInsufficient, err)
}

func TestService_ReserveCancel_UpdateStockError(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()

	items := []*model.Item{
		{SKU: 1, Count: 5},
	}

	repoMock := NewRepositoryMock(mc)

	repoMock.GetStockMock.
		When(ctx, model.SKUType(1)).
		Then(&model.Stock{SKU: 1, TotalCount: 20, ReservedCount: 10}, nil)

	expectedUpdateStocks := map[model.SKUType]*model.Stock{
		1: {SKU: 1, TotalCount: 20, ReservedCount: 5}, // 10 - 5
	}
	repoMock.UpdateStockMock.Set(func(ctx context.Context, stocks map[model.SKUType]*model.Stock) error {
		assert.True(t, compareStocks(expectedUpdateStocks, stocks))
		return errors.New("update error")
	})

	service := stockservice.NewService(repoMock)

	err := service.ReserveCancel(ctx, items)
	assert.Error(t, err)
	assert.EqualError(t, err, "update error")
}
