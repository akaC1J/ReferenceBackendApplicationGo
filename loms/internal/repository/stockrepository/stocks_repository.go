package stockrepository

import (
	"context"
	"fmt"
	"sync"

	appErr "route256/loms/internal/errors"
	"route256/loms/internal/model"
	"route256/loms/internal/service/stockservice"
)

var _ stockservice.Repository = (*Repository)(nil)

type Repository struct {
	stocks []*model.Stock
	mx     sync.Mutex
}

func NewEmptyRepository() *Repository {
	return &Repository{
		stocks: []*model.Stock{},
	}
}

func NewRepository(stocks []*model.Stock) *Repository {
	return &Repository{
		stocks: stocks,
	}
}

func (r *Repository) GetStock(_ context.Context, sku model.SKUType) (model.Stock, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	for _, stock := range r.stocks {
		if stock.SKU == sku {
			return *stock, nil
		}
	}
	return model.Stock{}, fmt.Errorf("stock with SKU %v: %w", sku, appErr.ErrNotFound)
}

func (r *Repository) UpdateStock(_ context.Context, stocks []model.Stock) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	existingStocks := make(map[model.SKUType]int)
	for indxStock, stock := range r.stocks {
		existingStocks[stock.SKU] = indxStock
	}

	for _, updateStock := range stocks {
		indxStock, found := existingStocks[updateStock.SKU]
		if !found {
			return fmt.Errorf("stock with SKU %v: %w", updateStock.SKU, appErr.ErrNotFound)
		}
		r.stocks[indxStock].TotalCount = updateStock.TotalCount
		r.stocks[indxStock].ReservedCount = updateStock.ReservedCount
	}

	return nil
}
