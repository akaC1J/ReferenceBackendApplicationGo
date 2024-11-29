-- +goose Up
-- +goose StatementBegin

CREATE TABLE orders
(
  id      BIGINT PRIMARY KEY,
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

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items, orders CASCADE;
-- +goose StatementEnd
