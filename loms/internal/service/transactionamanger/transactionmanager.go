package transactionmanager

import (
	"context"
	"github.com/jackc/pgx/v5"
	"route256/loms/internal/infra/database"
)

type ConnectionPooler interface {
	PickDefaultShard(ctx context.Context, readOnlyOperation bool) (*database.FallbackConnection, error)
	PickAllShards(ctx context.Context, readOnlyOperation bool) ([]*database.FallbackConnection, error)
}

type TransactionManager struct {
	pool ConnectionPooler
}

func NewTransactionManager(pool ConnectionPooler) *TransactionManager {
	return &TransactionManager{
		pool: pool,
	}
}

// Begin - начинает транзакцию и возвращает транзакцию и функцию закрытия
func (tm *TransactionManager) Begin(ctx context.Context) (pgx.Tx, error) {
	conn, err := tm.pool.PickDefaultShard(ctx, false)
	if err != nil {
		return nil, err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		conn.Release()
		return nil, err
	}

	return tx, nil
}

// BeginWithShard - начинает транзакцию и возвращает транзакцию и функцию закрытия на определенном шарде
func (tm *TransactionManager) BeginTransactionsOnAllShards(ctx context.Context) ([]pgx.Tx, error) {
	connections, err := tm.pool.PickAllShards(ctx, false)
	if err != nil {
		return nil, err
	}
	var transcations []pgx.Tx

	for _, conn := range connections {
		tx, err := conn.Begin(ctx)
		if err != nil {
			conn.Release()
			return nil, err
		}
		transcations = append(transcations, tx)
	}
	return transcations, nil
}
