-- name: CreatePurchase :one
INSERT INTO purchases(id, created_at, updated_at, chicken, price_per_chicken, user_id, farmer_id)
VALUES($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetPurchases :many

SELECT * FROM purchases;

-- name: GetPurchaseByID :one

SELECT * FROM purchases
WHERE id = $1;

-- name: DeletePurchase :exec

DELETE FROM purchases WHERE id = $1;

