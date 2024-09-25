-- orders.sql

-- name: SaveOrder :one
INSERT INTO orders (state, user_id)
VALUES ($1, $2)
RETURNING id, state, user_id;

-- name: UpdateOrder :one
UPDATE orders
SET state = $1
WHERE id = $2
RETURNING id, state, user_id;

-- name: GetOrderById :one
SELECT id, state, user_id
FROM orders
WHERE id = $1;
