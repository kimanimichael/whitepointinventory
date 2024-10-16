-- +goose Up

ALTER TABLE users
ADD COLUMN email TEXT UNIQUE NOT NULL
DEFAULT (md5(random()::text));

-- +goose Down

ALTER TABLE users DROP COLUMN email;