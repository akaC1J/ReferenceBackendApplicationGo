-- +goose Up
-- +goose StatementBegin
CREATE TABLE stock
(
  sku            BIGINT PRIMARY KEY,
  total_count    BIGINT NOT NULL,
  reserved_count BIGINT NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stock CASCADE;
-- +goose StatementEnd
