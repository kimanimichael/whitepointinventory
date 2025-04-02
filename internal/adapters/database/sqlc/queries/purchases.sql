-- name: CreatePurchase :one
INSERT INTO purchases(id, created_at, updated_at, chicken, price_per_chicken, user_id, farmer_id)
VALUES($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetPurchases :many

SELECT * FROM purchases ORDER BY created_at DESC;

-- name: GetPagedPurchases :many
SELECT * FROM purchases ORDER BY created_at DESC
OFFSET $1 LIMIT $2;

-- name: GetPurchasesCount :one
SELECT COUNT(*) AS total FROM purchases;


-- name: GetPurchaseByID :one

SELECT * FROM purchases
WHERE id = $1;

-- name: DeletePurchase :exec

DELETE FROM purchases WHERE id = $1;

-- name: GetMostRecentPurchase :one
SELECT * FROM purchases
ORDER BY created_at DESC
    LIMIT 1;

-- name: ChangePurchaseDate :exec
UPDATE purchases
SET updated_at = $2, created_at = $2
where purchases.id = $1;

