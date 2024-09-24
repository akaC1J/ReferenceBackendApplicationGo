package test

import (
	"context"
	"errors"
	appErrors "route256/loms/internal/errors"
	"route256/loms/internal/service/stockservice"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"route256/loms/internal/model"
)

func TestService_Reserve_Success(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()

	items := []*model.Item{
		{SKU: 1, Count: 10},
		{SKU: 2, Count: 5},
	}

	repoMock := NewRepositoryMock(mc)

	repoMock.GetStockMock.
		When(ctx, model.SKUType(1)).
		Then(&model.Stock{SKU: 1, TotalCount: 20, ReservedCount: 5}, nil)

	repoMock.GetStockMock.
		When(ctx, model.SKUType(2)).
		Then(&model.Stock{SKU: 2, TotalCount: 10, ReservedCount: 2}, nil)

	expectedUpdateStocks := map[model.SKUType]*model.Stock{
		1: {SKU: 1, TotalCount: 20, ReservedCount: 15}, // 5 + 10
		2: {SKU: 2, TotalCount: 10, ReservedCount: 7},  // 2 + 5
	}
	repoMock.UpdateStockMock.Set(func(ctx context.Context, stocks map[model.SKUType]*model.Stock) error {
		assert.True(t, compareStocks(expectedUpdateStocks, stocks))
		return nil
	})

	service := stockservice.NewService(repoMock)

	err := service.Reserve(ctx, items)
	assert.NoError(t, err)

	assert.Equal(t, uint64(2), repoMock.GetStockAfterCounter())
	assert.Equal(t, uint64(1), repoMock.UpdateStockAfterCounter())
}

func TestService_Reserve_NotEnoughStock(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()

	items := []*model.Item{
		{SKU: 1, Count: 15},
	}

	repoMock := NewRepositoryMock(mc)

	repoMock.GetStockMock.
		When(ctx, model.SKUType(1)).
		Then(&model.Stock{SKU: 1, TotalCount: 10, ReservedCount: 0}, nil)

	service := stockservice.NewService(repoMock)

	err := service.Reserve(ctx, items)
	assert.Error(t, err)
	assert.ErrorIs(t, err, appErrors.ErrStockInsufficient)

	assert.Equal(t, uint64(1), repoMock.GetStockAfterCounter())
	assert.Equal(t, uint64(0), repoMock.UpdateStockAfterCounter())
}

func TestService_Reserve_GetStockError(t *testing.T) {
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

	err := service.Reserve(ctx, items)
	assert.Error(t, err)
	assert.EqualError(t, err, "database error")

	assert.Equal(t, uint64(1), repoMock.GetStockAfterCounter())
	assert.Equal(t, uint64(0), repoMock.UpdateStockAfterCounter())
}

func TestService_Reserve_UpdateStockError(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()

	items := []*model.Item{
		{SKU: 1, Count: 5},
	}

	repoMock := NewRepositoryMock(mc)

	repoMock.GetStockMock.
		When(ctx, model.SKUType(1)).
		Then(&model.Stock{SKU: 1, TotalCount: 10, ReservedCount: 2}, nil)

	expectedUpdateStocks := map[model.SKUType]*model.Stock{
		1: {SKU: 1, TotalCount: 10, ReservedCount: 7}, // 2 + 5
	}

	repoMock.UpdateStockMock.Set(func(ctx context.Context, stocks map[model.SKUType]*model.Stock) error {
		assert.True(t, compareStocks(expectedUpdateStocks, stocks))
		return errors.New("update error")
	})
	service := stockservice.NewService(repoMock)

	err := service.Reserve(ctx, items)
	assert.Error(t, err)
	assert.EqualError(t, err, "update error")

	assert.Equal(t, uint64(1), repoMock.GetStockAfterCounter())
	assert.Equal(t, uint64(1), repoMock.UpdateStockAfterCounter())
}

func TestService_Reserve_DuplicateSKUs(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()

	items := []*model.Item{
		{SKU: 1, Count: 5},
		{SKU: 1, Count: 3},
	}

	repoMock := NewRepositoryMock(mc)

	repoMock.GetStockMock.
		When(ctx, model.SKUType(1)).
		Then(&model.Stock{SKU: 1, TotalCount: 10, ReservedCount: 2}, nil)

	expectedUpdateStocks := map[model.SKUType]*model.Stock{
		1: {SKU: 1, TotalCount: 10, ReservedCount: 10}, // 2 + (5 + 3)
	}
	repoMock.UpdateStockMock.Set(func(ctx context.Context, stocks map[model.SKUType]*model.Stock) error {
		assert.True(t, compareStocks(expectedUpdateStocks, stocks))
		return nil
	})

	service := stockservice.NewService(repoMock)

	err := service.Reserve(ctx, items)
	assert.NoError(t, err)

	assert.Equal(t, uint64(1), repoMock.GetStockAfterCounter())
	assert.Equal(t, uint64(1), repoMock.UpdateStockAfterCounter())
}
