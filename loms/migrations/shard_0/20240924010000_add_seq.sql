-- +goose Up
-- +goose StatementBegin

CREATE SEQUENCE order_id_manual_seq_shard_0 INCREMENT 1000 START 1000; -- 1000 > number of buckets

ALTER TABLE orders ALTER COLUMN id SET DEFAULT nextval('order_id_manual_seq_shard_0');

ALTER SEQUENCE order_id_manual_seq_shard_0 OWNED BY orders.id;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP SEQUENCE order_id_manual_seq_shard_0;

ALTER TABLE orders ALTER COLUMN id DROP DEFAULT;

-- +goose StatementEnd
