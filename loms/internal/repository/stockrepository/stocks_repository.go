package stockrepository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go.uber.org/multierr"
	apperrors "route256/loms/internal/errors"
	"route256/loms/internal/infra/database"
	"route256/loms/internal/model"
	"route256/loms/internal/repository"
)

type Repository struct {
	pool ConnectionPooler
}

type ConnectionPooler interface {
	PickDefaultShard(ctx context.Context, readOnly bool) (*database.FallbackConnection, error)
}

func NewRepository(pool ConnectionPooler) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) GetStocks(ctx context.Context, sku []model.SKUType) (stocks []*model.Stock, err error) {
	intSku := make([]int64, 0, len(sku))
	for _, s := range sku {
		intSku = append(intSku, int64(s))
	}
	conn, err := r.pool.PickDefaultShard(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("unable to acquire a connection: %w", err)
	}
	repositoryStocks, err := New(conn).GetStockBySkus(ctx, intSku)
	if err != nil {
		return nil, err
	}
	if len(repositoryStocks) != len(sku) {
		unfindedSku := make([]model.SKUType, 0)
		for _, s := range sku {
			find := false
			for _, rs := range repositoryStocks {
				if int64(s) == rs.Sku {
					find = true
					break
				}
			}
			if !find {
				unfindedSku = append(unfindedSku, s)
			}
		}
		return nil, fmt.Errorf("sku %v in database", apperrors.ErrNotFound)
	}
	stocks, err = repackRepositoryStockToModelStock(repositoryStocks)
	if err != nil {
		return nil, err
	}
	return stocks, nil
}

func (r *Repository) UpdateStock(ctx context.Context, stocks map[model.SKUType]*model.Stock) error {
	conn, err := r.pool.PickDefaultShard(ctx, false)
	if err != nil {
		return fmt.Errorf("unable to acquire a connection: %w", err)
	}
	defer conn.Release()
	err = pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) (err error) {
		repTx := New(tx)
		updateParamStock := repackStocksMapToUpdateStockParam(stocks)
		result, err := repTx.db.Exec(ctx, updateStockInfo, updateParamStock.Skus,
			updateParamStock.TotalCounts,
			updateParamStock.ReservedCounts)
		if err != nil {
			return err
		}
		if result.RowsAffected() != int64(len(stocks)) {
			return fmt.Errorf("expected %v rows affected, got %v", len(stocks), result.RowsAffected())
		}
		return nil
	})
	return err
}

func repackStocksMapToUpdateStockParam(stocks map[model.SKUType]*model.Stock) *UpdateStockInfoParams {
	skus := make([]int64, 0, len(stocks))
	totalCounts := make([]int64, 0, len(stocks))
	reservedCounts := make([]int64, 0, len(stocks))
	//safe cast uint32 to int64
	for _, stock := range stocks {
		skus = append(skus, int64(stock.SKU))
		totalCounts = append(totalCounts, int64(stock.TotalCount))
		reservedCounts = append(reservedCounts, int64(stock.ReservedCount))
	}
	return &UpdateStockInfoParams{
		Skus:           skus,
		TotalCounts:    totalCounts,
		ReservedCounts: reservedCounts,
	}

}
func repackRepositoryStockToModelStock(stocksFromDB []*Stock) ([]*model.Stock, error) {
	var stocks []*model.Stock
	for _, s := range stocksFromDB {
		var resultErr error
		safeSku, err := repository.SafeInt64ToUint32(s.Sku)
		resultErr = multierr.Append(resultErr, err)
		safeTotalCount, err := repository.SafeInt64ToUint32(s.TotalCount)
		resultErr = multierr.Append(resultErr, err)
		safeReservedCount, err := repository.SafeInt64ToUint32(s.ReservedCount)
		resultErr = multierr.Append(resultErr, err)
		if resultErr != nil {
			return nil, resultErr
		}
		stocks = append(stocks, &model.Stock{
			SKU:           model.SKUType(safeSku),
			TotalCount:    safeTotalCount,
			ReservedCount: safeReservedCount,
		})
	}
	return stocks, nil
}
