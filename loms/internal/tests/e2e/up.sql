-- up.sql

CREATE TYPE ORDER_STATUS AS ENUM ('NEW', 'AWAITING_PAYMENT', 'FAILED', 'PAYED', 'CANCELLED');

-- Create tables
CREATE TABLE orders
(
  id      BIGSERIAL PRIMARY KEY,
  state   ORDER_STATUS NOT NULL,
  user_id BIGINT       NOT NULL
);

CREATE TABLE items
(
  sku      BIGINT NOT NULL,
  order_id BIGINT NOT NULL,
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

INSERT INTO stock (sku, total_count, reserved_count)
VALUES (1, 150, 10),
       (2, 200, 20),
       (3, 250, 30),
       (4, 300, 40),
       (5, 350, 50);
