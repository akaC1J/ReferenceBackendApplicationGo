-- +goose Up
-- +goose StatementBegin
CREATE TYPE ORDER_STATUS AS ENUM ('NEW', 'AWAITING PAYMENT', 'FAILED', 'PAYED', 'CANCELLED');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE ORDER_STATUS
-- +goose StatementEnd
