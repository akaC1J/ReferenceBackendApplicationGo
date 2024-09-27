-- name: UpdateStockInfo :exec
UPDATE stock
SET total_count    = data.total_count,
    reserved_count = data.reserved_count
FROM (SELECT unnest(@skus::bigint[])             AS sku,
             unnest(@total_counts::bigint[])     AS total_count,
             unnest(@reserved_counts:: bigint[]) AS reserved_count) AS data
WHERE stock.sku = data.sku;

-- name: GetStockBySkus :many
SELECT sku,
       total_count,
       reserved_count
FROM stock
WHERE sku = ANY(@skus :: bigint[]);
