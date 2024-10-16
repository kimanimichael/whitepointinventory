-- +goose Up

CREATE TABLE farmers(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT UNIQUE NOT NULL,
    chicken_balance INT,
    cash_balance INT
);

-- +goose Down

DROP TABLE farmers;