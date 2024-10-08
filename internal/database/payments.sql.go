// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: payments.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createPayment = `-- name: CreatePayment :one
INSERT INTO payments (id, created_at, updated_at, cash_paid, price_per_chicken_paid, user_id, farmer_id)
VALUES($1, $2, $3, $4, $5, $6, $7)
RETURNING id, created_at, updated_at, cash_paid, price_per_chicken_paid, user_id, farmer_id
`

type CreatePaymentParams struct {
	ID                  uuid.UUID
	CreatedAt           time.Time
	UpdatedAt           time.Time
	CashPaid            int32
	PricePerChickenPaid int32
	UserID              uuid.UUID
	FarmerID            uuid.UUID
}

func (q *Queries) CreatePayment(ctx context.Context, arg CreatePaymentParams) (Payment, error) {
	row := q.db.QueryRowContext(ctx, createPayment,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.CashPaid,
		arg.PricePerChickenPaid,
		arg.UserID,
		arg.FarmerID,
	)
	var i Payment
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CashPaid,
		&i.PricePerChickenPaid,
		&i.UserID,
		&i.FarmerID,
	)
	return i, err
}

const deletePayments = `-- name: DeletePayments :exec
DELETE FROM payments WHERE id = $1
`

func (q *Queries) DeletePayments(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deletePayments, id)
	return err
}

const getMostRecentPayment = `-- name: GetMostRecentPayment :one
SELECT id, created_at, updated_at, cash_paid, price_per_chicken_paid, user_id, farmer_id FROM payments
ORDER BY created_at DESC
LIMIT 1
`

func (q *Queries) GetMostRecentPayment(ctx context.Context) (Payment, error) {
	row := q.db.QueryRowContext(ctx, getMostRecentPayment)
	var i Payment
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CashPaid,
		&i.PricePerChickenPaid,
		&i.UserID,
		&i.FarmerID,
	)
	return i, err
}

const getPaymentByID = `-- name: GetPaymentByID :one
SELECT id, created_at, updated_at, cash_paid, price_per_chicken_paid, user_id, farmer_id FROM payments
WHERE id = $1
`

func (q *Queries) GetPaymentByID(ctx context.Context, id uuid.UUID) (Payment, error) {
	row := q.db.QueryRowContext(ctx, getPaymentByID, id)
	var i Payment
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CashPaid,
		&i.PricePerChickenPaid,
		&i.UserID,
		&i.FarmerID,
	)
	return i, err
}

const getPayments = `-- name: GetPayments :many
SELECT id, created_at, updated_at, cash_paid, price_per_chicken_paid, user_id, farmer_id FROM payments ORDER BY created_at DESC
`

func (q *Queries) GetPayments(ctx context.Context) ([]Payment, error) {
	rows, err := q.db.QueryContext(ctx, getPayments)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Payment
	for rows.Next() {
		var i Payment
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.CashPaid,
			&i.PricePerChickenPaid,
			&i.UserID,
			&i.FarmerID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
