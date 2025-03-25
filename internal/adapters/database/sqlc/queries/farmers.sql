-- name: CreateFarmer :one
INSERT INTO farmers(id, created_at, updated_at, name, chicken_balance, cash_balance)
VALUES($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFarmerByName :one
SELECT * FROM farmers where name = $1;

-- name: DeleteFarmers :exec
DELETE FROM farmers where id = $1;

-- name: IncreaseChickenOwed :exec
UPDATE farmers
SET chicken_balance = COALESCE(chicken_balance, 0) + $1
WHERE farmers.id = $2;

-- name: IncreaseCashOwed :exec
UPDATE farmers
SET cash_balance = COALESCE(cash_balance, 0) + ($1)
WHERE farmers.id = $2;

-- name: DecreaseChickenOwed :exec
UPDATE farmers
SET chicken_balance = COALESCE(chicken_balance, 0) - $1
WHERE farmers.id = $2;

-- name: DecreaseCashOwed :exec
UPDATE farmers
SET cash_balance = COALESCE(cash_balance, 0) - ($1)
WHERE farmers.id = $2;

-- name: GetFarmerByID :one
SELECT * FROM farmers where id = $1;

-- name: GetFarmers :many
SELECT * FROM farmers ORDER BY updated_at DESC;

-- name: GetPagedFarmers :many
SELECT * FROM farmers OFFSET $1 LIMIT $2;

-- name: SetFarmerBalances :exec
UPDATE farmers
SET chicken_balance = $2, cash_balance = $3
where farmers.name = $1;

-- name: GetFarmerCount :one
SELECT COUNT(*) AS total FROM farmers;

-- name: MarkFarmerAsUpdated :exec
UPDATE farmers
SET updated_at = NOW()
WHERE id = $1;