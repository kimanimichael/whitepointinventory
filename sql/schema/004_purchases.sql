-- +goose Up

CREATE TABLE purchases(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    chicken INT NOT NULL,
    price_per_chicken INT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    farmer_id UUID NOT NULL REFERENCES farmers(id)
);

-- +goose Down

DROP TABLE purchases;