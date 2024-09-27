package orderrepository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/multierr"
	appErr "route256/loms/internal/errors"
	"route256/loms/internal/infra"
	"route256/loms/internal/model"
	"route256/loms/internal/repository"
)

type Repository struct {
	q    *Queries
	pool ConnectionPooler
}

type ConnectionPooler interface {
	DBTX
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
}

func NewRepository(pool ConnectionPooler) *Repository {
	return &Repository{
		q:    New(pool),
		pool: pool,
	}
}

func (r *Repository) SaveOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	ctx = context.WithValue(ctx, infra.ReadOnlyKey, false)
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to acquire a connection: %w", err)
	}
	defer conn.Release()
	err = pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
		repTx := New(tx)
		orderID, err := repTx.SaveOrder(ctx, &SaveOrderParams{
			State:  OrderStatus(order.State()),
			UserID: order.UserId,
		})
		if err != nil {
			return err
		}
		err = repTx.SaveItems(ctx, repackItemsToSaveItemParams(order.Items, orderID))
		if err != nil {
			return err
		}
		order.ID = orderID
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("unable to save order: %w", err)
	}
	return order, nil
}

func (r *Repository) UpdateOrder(ctx context.Context, order *model.Order) error {
	ctx = context.WithValue(ctx, infra.ReadOnlyKey, false)
	_, err := r.q.UpdateOrder(ctx, &UpdateOrderParams{
		UserID:  order.UserId,
		State:   OrderStatus(order.State()),
		OrderID: order.ID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("order with ID %v: %w", order.ID, appErr.ErrNotFound)
		}
		return fmt.Errorf("unable to update order: %w", err)
	}
	return nil
}

func (r *Repository) GetById(ctx context.Context, orderID int64) (*model.Order, error) {
	ctx = context.WithValue(ctx, infra.ReadOnlyKey, true)
	orderFromDB, err := r.q.GetOrderById(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("unable to get order by ID: %w", err)
	}
	order, err := repackOrderFromDBToOrder(orderFromDB)
	if err != nil {
		return nil, fmt.Errorf("unable to repack order from DB: %w", err)
	}
	return order, nil
}

func repackOrderFromDBToOrder(orderFromDB []*GetOrderByIdRow) (*model.Order, error) {
	if len(orderFromDB) == 0 {
		return nil, appErr.ErrNotFound
	}

	order := &model.Order{
		ID:     orderFromDB[0].ID,
		UserId: orderFromDB[0].UserID,
	}
	err := order.SetState(model.StateType(orderFromDB[0].State))
	if err != nil {
		return nil, err
	}

	items := make([]*model.Item, 0, len(orderFromDB))
	for _, item := range orderFromDB {
		var resultError error
		safeSku, err := repository.SafeInt64ToUint32(item.Sku)
		resultError = multierr.Append(resultError, err)
		safeCount, err := repository.SafeInt64ToUint32(item.Count)
		resultError = multierr.Append(resultError, err)
		if resultError != nil {
			return nil, resultError
		}
		items = append(items, &model.Item{
			SKU:   model.SKUType(safeSku),
			Count: safeCount,
		})
	}
	order.Items = items
	return order, nil
}

func repackItemsToSaveItemParams(items []*model.Item, orderID int64) *SaveItemsParams {
	skus := make([]int64, 0, len(items))
	counts := make([]int64, 0, len(items))
	for _, item := range items {
		//unsafe conversion
		skus = append(skus, int64(item.SKU))
		counts = append(counts, int64(item.Count))
	}
	return &SaveItemsParams{
		Skus:    skus,
		Counts:  counts,
		OrderID: orderID,
	}

}
