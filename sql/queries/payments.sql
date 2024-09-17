-- name: CreatePayment :one
INSERT INTO payments (id, created_at, updated_at, cash_paid, price_per_chicken_paid, user_id, farmer_id)
VALUES($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetPayments :many
SELECT * FROM payments ORDER BY created_at DESC;

-- name: GetPaymentByID :one
SELECT * FROM payments
WHERE id = $1;

-- name: DeletePayments :exec
DELETE FROM payments WHERE id = $1;

-- name: GetMostRecentPayment :one
SELECT * FROM payments
ORDER BY created_at DESC
LIMIT 1;