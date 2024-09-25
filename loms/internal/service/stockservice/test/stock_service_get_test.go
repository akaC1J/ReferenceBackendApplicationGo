package test

import (
	"context"
	"errors"
	appErors "route256/loms/internal/errors"
	"route256/loms/internal/model"
	"route256/loms/internal/service/stockservice"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestService_GetBySKUAvailableCount_Success(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()
	sku := model.SKUType(1)

	repoMock := NewRepositoryMock(mc)
	repoMock.GetStockMock.
		When(ctx, sku).
		Then(&model.Stock{SKU: sku, TotalCount: 20, ReservedCount: 5}, nil)

	service := stockservice.NewService(repoMock)

	availableCount, err := service.GetBySKUAvailableCount(ctx, sku)
	assert.NoError(t, err)
	assert.Equal(t, uint64(15), availableCount) // 20 - 5
}

func TestService_GetBySKUAvailableCount_GetStockError(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()
	sku := model.SKUType(1)

	repoMock := NewRepositoryMock(mc)
	repoMock.GetStockMock.
		When(ctx, sku).
		Then(&model.Stock{}, errors.New("database error"))

	service := stockservice.NewService(repoMock)

	availableCount, err := service.GetBySKUAvailableCount(ctx, sku)
	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
	assert.Equal(t, uint64(0), availableCount)
}

func TestService_GetBySKUAvailableCount_NegativeAvailableCount(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()
	sku := model.SKUType(1)

	repoMock := NewRepositoryMock(mc)
	repoMock.GetStockMock.
		When(ctx, sku).
		Then(&model.Stock{SKU: sku, TotalCount: 5, ReservedCount: 10}, nil)

	service := stockservice.NewService(repoMock)

	availableCount, err := service.GetBySKUAvailableCount(ctx, sku)
	assert.Error(t, err)
	assert.IsType(t, appErors.ErrStockInsufficient, err)
	assert.Equal(t, uint64(0), availableCount)
}
