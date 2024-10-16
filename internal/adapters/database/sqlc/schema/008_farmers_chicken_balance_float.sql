-- +goose Up

ALTER TABLE farmers
ALTER COLUMN chicken_balance TYPE FLOAT USING chicken_balance::FLOAT;

-- +goose Down
ALTER TABLE farmers
ALTER COLUMN chicken_balance TYPE INT USING chicken_balance::INT;