-- name: SaveOutboxEvent :exec
INSERT INTO outbox (order_id, payload)
VALUES ($1, $2);

-- name: GetPendingOutboxEvents :many
SELECT id, order_id, payload, created_at
FROM outbox
WHERE processed = FALSE
ORDER BY created_at
LIMIT $1;

-- name: MarkOutboxEventProcessed :exec
UPDATE outbox
SET processed = TRUE
WHERE id = $1;
