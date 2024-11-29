package transactionmanager

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"route256/loms/internal/infra/database"
)

type ConnectionPooler interface {
	PickDefaultShard(ctx context.Context, readOnlyOperation bool) (*database.FallbackConnection, error)
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
func (tm *TransactionManager) Begin(ctx context.Context) (pgx.Tx, func(), error) {
	conn, err := tm.pool.PickDefaultShard(ctx, false)
	if err != nil {
		return nil, nil, err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		conn.Release()
		return nil, nil, err
	}

	// Закрывающая функция для безопасного отката и освобождения соединения
	closeFunc := func() {
		if err := tx.Rollback(context.Background()); err != nil && err != pgx.ErrTxClosed {
			log.Printf("Failed to rollback transaction: %v", err)
		}
		conn.Release()
	}

	return tx, closeFunc, nil
}

// Commit - выполняет коммит и освобождает соединение
func (tm *TransactionManager) Commit(tx pgx.Tx, closeFunc func()) error {
	defer closeFunc() // Всегда освобождает ресурсы

	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}
