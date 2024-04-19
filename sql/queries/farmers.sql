-- name: CreateFarmer :one
INSERT INTO farmers(id, created_at, updated_at, name, chicken_balance, cash_balance)
VALUES($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFarmerByName :one
SELECT * FROM farmers where name = $1;
