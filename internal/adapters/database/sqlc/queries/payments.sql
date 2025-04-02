-- name: CreatePayment :one
INSERT INTO payments (id, created_at, updated_at, cash_paid, price_per_chicken_paid, user_id, farmer_id)
VALUES($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetPayments :many
SELECT * FROM payments ORDER BY created_at DESC;

-- name: GetPagedPayments :many
SELECT * FROM payments ORDER BY created_at DESC
OFFSET $1 LIMIT $2;

-- name: GetPaymentCount :one
SELECT COUNT(*) AS total FROM payments;

-- name: GetPaymentByID :one
SELECT * FROM payments
WHERE id = $1;

-- name: DeletePayments :exec
DELETE FROM payments WHERE id = $1;

-- name: GetMostRecentPayment :one
SELECT * FROM payments
ORDER BY created_at DESC
    LIMIT 1;

-- name: ChangePaymentDate :exec
UPDATE payments
SET updated_at = $2, created_at = $2
where payments.id = $1;