package orderrepository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go.uber.org/multierr"
	appErr "route256/loms/internal/errors"
	"route256/loms/internal/infra/database"
	"route256/loms/internal/model"
	"route256/loms/internal/repository"
	"route256/loms/internal/repository/outboxrepository"
	"slices"
)

type Repository struct {
	pool             ConnectionPooler
	outboxRepository *outboxrepository.Repository
}

type ConnectionPooler interface {
	PickConnFromUserId(ctx context.Context, userId int64, readOnlyOperation bool) (*database.FallbackConnection, error)
	PickConnFromOrderId(ctx context.Context, orderID int64, readOnlyOperation bool) (*database.FallbackConnection, error)
	PickAllShards(ctx context.Context, readOnlyOperation bool) ([]*database.FallbackConnection, error)
}

func NewRepository(pool ConnectionPooler, oR *outboxrepository.Repository) *Repository {
	return &Repository{
		pool:             pool,
		outboxRepository: oR,
	}
}

func (r *Repository) SaveOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	conn, err := r.pool.PickConnFromUserId(ctx, order.UserId, false)
	if err != nil {
		return nil, fmt.Errorf("unable to acquire a connection: %w", err)
	}
	defer conn.Release()

	payload, err := json.Marshal(order)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal order: %w", err)
	}
	err = pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
		repTx := New(tx)
		orderID, err := repTx.SaveOrder(ctx, &SaveOrderParams{
			State:  OrderStatus(order.State),
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

		err = r.outboxRepository.SaveOutboxEvent(ctx, tx, &model.OutboxEvent{
			OrderID: order.ID,
			Payload: string(payload),
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("unable to save order: %w", err)
	}
	return order, nil
}

func (r *Repository) UpdateOrder(ctx context.Context, order *model.Order) error {
	conn, err := r.pool.PickConnFromUserId(ctx, order.UserId, false)
	if err != nil {
		return fmt.Errorf("unable to acquire a connection: %w", err)
	}
	defer conn.Release()

	payload, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("unable to marshal order: %w", err)
	}

	err = pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
		repTx := New(tx)
		_, err = repTx.UpdateOrder(ctx, &UpdateOrderParams{
			UserID:  order.UserId,
			State:   OrderStatus(order.State),
			OrderID: order.ID,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("order with ID %v: %w", order.ID, appErr.ErrNotFound)
			}
			return fmt.Errorf("unable to update order: %w", err)
		}
		err = r.outboxRepository.SaveOutboxEvent(ctx, tx, &model.OutboxEvent{
			OrderID: order.ID,
			Payload: string(payload),
		})
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (r *Repository) GetById(ctx context.Context, orderID int64) (*model.Order, error) {
	conn, err := r.pool.PickConnFromOrderId(ctx, orderID, false)
	orderFromDB, err := New(conn).GetOrderById(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("unable to get order by ID: %w", err)
	}
	order, err := repackOrderFromDBToOrder(orderFromDB)
	if err != nil {
		return nil, fmt.Errorf("unable to repack order from DB: %w", err)
	}
	return order, nil
}

func (r *Repository) GetAllOrders(ctx context.Context) ([]*model.Order, error) {
	connections, err := r.pool.PickAllShards(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("error pick all shards: %w", err)
	}
	orders := make([]*model.Order, 0)
	for id, conn := range connections {
		ordersFromDb, err := New(conn).GetAllOrders(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to get orders from shard id %d: %w", id, err)
		}
		toOrders, err := repackOrdersFromDBToOrders(ordersFromDb)
		if err != nil {
			return nil, fmt.Errorf("unable to repack order from DB for shard id %d: %w", id, err)
		}
		orders = append(orders, toOrders...)
	}

	slices.SortFunc(orders, func(i, j *model.Order) int {
		return -int(i.ID - j.ID)
	})

	return orders, nil
}

func repackOrdersFromDBToOrders(ordersFromDB []*GetAllOrdersRow) ([]*model.Order, error) {
	ordersRes := make([]*model.Order, 0)
	orders := make(map[int64][]*GetOrderByIdRow)
	for _, orderFromDB := range ordersFromDB {
		orders[orderFromDB.ID] = append(orders[orderFromDB.ID], &GetOrderByIdRow{
			ID:     orderFromDB.ID,
			State:  orderFromDB.State,
			UserID: orderFromDB.UserID,
			Sku:    orderFromDB.Sku,
			Count:  orderFromDB.Count,
		})
	}

	for key := range orders {
		order, err := repackOrderFromDBToOrder(orders[key])
		if err != nil {
			return nil, err
		}
		ordersRes = append(ordersRes, order)
	}

	return ordersRes, nil
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
