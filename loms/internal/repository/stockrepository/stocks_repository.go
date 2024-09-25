package stockrepository

import (
	"context"
	"fmt"
	"sync"

	appErr "route256/loms/internal/errors"
	"route256/loms/internal/model"
)

type Repository struct {
	stocks map[model.SKUType]model.Stock
	mx     sync.Mutex
}

func NewRepository(stocks map[model.SKUType]model.Stock) *Repository {
	return &Repository{
		stocks: stocks,
	}
}

func (r *Repository) GetStock(_ context.Context, sku model.SKUType) (*model.Stock, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	stock, found := r.stocks[sku]
	if !found {
		return nil, fmt.Errorf("stock with SKU %v: %w", sku, appErr.ErrNotFound)
	}
	return &stock, nil
}

func (r *Repository) UpdateStock(_ context.Context, stocks map[model.SKUType]*model.Stock) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	for sku, updateStock := range stocks {
		stock, found := r.stocks[sku]
		if !found {
			return fmt.Errorf("stock with SKU %v: %w", sku, appErr.ErrNotFound)
		}
		stock.TotalCount = updateStock.TotalCount
		stock.ReservedCount = updateStock.ReservedCount
		r.stocks[sku] = stock
	}

	return nil
}
