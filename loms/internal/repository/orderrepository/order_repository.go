package orderrepository

import (
	"context"
	"fmt"
	"sync"

	appErr "route256/loms/internal/errors"
	"route256/loms/internal/model"
	"route256/loms/internal/service/orderservice"
)

var _ orderservice.Repository = (*Repository)(nil)

type Repository struct {
	orders []model.Order
	mx     sync.Mutex
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) SaveOrder(_ context.Context, order *model.Order) (*model.Order, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	order.ID = int64(len(r.orders) + 1)
	r.orders = append(r.orders, *order)
	return order, nil
}

func (r *Repository) UpdateOrder(_ context.Context, order *model.Order) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	for i, o := range r.orders {
		if o.ID == order.ID {
			r.orders[i] = *order
			return nil
		}
	}
	return fmt.Errorf("order with ID %v: %w", order.ID, appErr.ErrNotFound)
}

func (r *Repository) GetById(_ context.Context, orderID int64) (*model.Order, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	for _, order := range r.orders {
		if order.ID == orderID {
			return &order, nil
		}
	}
	return nil, fmt.Errorf("order with ID %v: %w", orderID, appErr.ErrNotFound)
}
