package orderrepository

import (
	"context"
	"fmt"
	"sync"

	appErr "route256/loms/internal/errors"
	"route256/loms/internal/model"
)

type Repository struct {
	orders map[int64]model.Order
	mx     sync.Mutex
}

func NewRepository() *Repository {
	return &Repository{orders: make(map[int64]model.Order)}
}

func (r *Repository) SaveOrder(_ context.Context, order *model.Order) (*model.Order, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	order.ID = int64(len(r.orders) + 1)
	r.orders[order.ID] = *order
	return order, nil
}

func (r *Repository) UpdateOrder(_ context.Context, order *model.Order) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	_, ok := r.orders[order.ID]
	if !ok {
		return fmt.Errorf("order with ID %v: %w", order.ID, appErr.ErrNotFound)
	}
	r.orders[order.ID] = *order
	return nil
}

func (r *Repository) GetById(_ context.Context, orderID int64) (*model.Order, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	order, ok := r.orders[orderID]
	if ok {
		return &order, nil
	}
	return nil, fmt.Errorf("order with ID %v: %w", orderID, appErr.ErrNotFound)
}
