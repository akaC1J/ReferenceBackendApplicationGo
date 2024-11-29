-- +goose Up
-- +goose StatementBegin

CREATE SEQUENCE order_id_manual_seq_shard_1 INCREMENT 1000 START 1001; -- 1000 > number of buckets

ALTER TABLE orders ALTER COLUMN id SET DEFAULT nextval('order_id_manual_seq_shard_1');

ALTER SEQUENCE order_id_manual_seq_shard_1 OWNED BY orders.id;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP SEQUENCE order_id_manual_seq_shard_1;

ALTER TABLE orders ALTER COLUMN id DROP DEFAULT;

-- +goose StatementEnd
