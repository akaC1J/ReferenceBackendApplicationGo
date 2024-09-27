-- down.sql

-- Delete data from stock table
DELETE FROM stock WHERE sku IN (773297411, 1002, 1003, 1004, 1005);

-- Drop tables
DROP TABLE IF EXISTS items, orders, stock CASCADE;

-- Drop custom type
DROP TYPE IF EXISTS ORDER_STATUS;
