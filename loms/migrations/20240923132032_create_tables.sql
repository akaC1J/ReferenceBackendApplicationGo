-- +goose Up
-- +goose StatementBegin

CREATE TABLE orders
(
  id      BIGSERIAL PRIMARY KEY,
  state   ORDER_STATUS NOT NULL,
  user_id BIGINT          NOT NULL
);

CREATE TABLE items
(
  sku      BIGINT NOT NULL ,
  order_id BIGINT NOT NULL ,
  count    BIGINT NOT NULL,
  PRIMARY KEY (order_id, sku),
  FOREIGN KEY (order_id) REFERENCES orders (id)
);

CREATE TABLE stock
(
  sku            BIGINT PRIMARY KEY,
  total_count    BIGINT NOT NULL,
  reserved_count BIGINT NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items, orders, stock CASCADE;
-- +goose StatementEnd
