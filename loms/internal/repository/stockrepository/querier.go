// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package stockrepository

import (
	"context"
)

type Querier interface {
	GetStockBySkus(ctx context.Context, skus []int64) ([]*Stock, error)
	UpdateStockInfo(ctx context.Context, arg *UpdateStockInfoParams) error
}

var _ Querier = (*Queries)(nil)
