-- name: SaveOrder :one
INSERT INTO orders (state, user_id)
VALUES (@state, @user_id)
RETURNING id;

-- name: SaveItems :exec
INSERT INTO items (sku, count, order_id)
SELECT unnest(@skus::bigint[]), unnest(@counts::bigint[]), @order_id;

-- name: UpdateOrder :one
UPDATE orders
SET state   = @state,
    user_id = @user_id
WHERE id = @order_id
RETURNING id;

-- name: GetOrderById :many
SELECT orders.id,
       orders.state,
       orders.user_id,
       i.sku,
       i.count
FROM orders
JOIN items i on orders.id = i.order_id
WHERE orders.id = @order_id;
