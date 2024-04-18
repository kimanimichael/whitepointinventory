-- +goose Up

ALTER TABLE users ADD CONSTRAINT unique_name UNIQUE(name);

-- +goose Down

ALTER TABLE users DROP CONSTRAINT unique_name;