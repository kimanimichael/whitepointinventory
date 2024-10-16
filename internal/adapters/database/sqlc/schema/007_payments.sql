-- +goose Up

CREATE TABLE payments(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    cash_paid INT NOT NULL,
    price_per_chicken_paid INT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    farmer_id UUID NOT NULL REFERENCES farmers(id)
);

-- +goose Down

DROP TABLE payments;