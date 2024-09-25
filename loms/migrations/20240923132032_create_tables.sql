-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
  id SERIAL PRIMARY KEY
  -- другие поля пользователя (например, имя, email и т.д.)
);

CREATE TABLE orders
(
  id      SERIAL PRIMARY KEY,
  state   ORDER_STATUS NOT NULL,
  user_id INT          NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE items
(
  id       SERIAL PRIMARY KEY,
  sku      INTEGER NOT NULL,
  count    INTEGER NOT NULL,
  order_id INT,
  FOREIGN KEY (order_id) REFERENCES orders (id)
);

CREATE TABLE stock
(
  sku            INTEGER PRIMARY KEY,
  total_count    INTEGER NOT NULL,
  reserved_count INTEGER NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items, orders, users, stock CASCADE;
-- +goose StatementEnd
