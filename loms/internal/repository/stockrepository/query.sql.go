// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package stockrepository

import (
	"context"
)

const getStockBySkus = `-- name: GetStockBySkus :many
SELECT sku,
       total_count,
       reserved_count
FROM stock
WHERE sku = ANY($1 :: bigint[])
`

func (q *Queries) GetStockBySkus(ctx context.Context, skus []int64) ([]*Stock, error) {
	rows, err := q.db.Query(ctx, getStockBySkus, skus)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Stock
	for rows.Next() {
		var i Stock
		if err := rows.Scan(&i.Sku, &i.TotalCount, &i.ReservedCount); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateStockInfo = `-- name: UpdateStockInfo :exec
UPDATE stock
SET total_count    = data.total_count,
    reserved_count = data.reserved_count
FROM (SELECT unnest($1::bigint[])             AS sku,
             unnest($2::bigint[])     AS total_count,
             unnest($3:: bigint[]) AS reserved_count) AS data
WHERE stock.sku = data.sku
`

type UpdateStockInfoParams struct {
	Skus           []int64
	TotalCounts    []int64
	ReservedCounts []int64
}

func (q *Queries) UpdateStockInfo(ctx context.Context, arg *UpdateStockInfoParams) error {
	_, err := q.db.Exec(ctx, updateStockInfo, arg.Skus, arg.TotalCounts, arg.ReservedCounts)
	return err
}
