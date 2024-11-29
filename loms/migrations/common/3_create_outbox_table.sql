-- +goose Up
-- +goose StatementBegin
CREATE TABLE outbox
(
  id         SERIAL PRIMARY KEY,
  order_id   BIGINT      NOT NULL,
  payload    TEXT        NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  processed  BOOLEAN                  DEFAULT FALSE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS outbox CASCADE;
-- +goose StatementEnd
