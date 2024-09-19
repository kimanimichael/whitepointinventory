// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: purchases.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createPurchase = `-- name: CreatePurchase :one
INSERT INTO purchases(id, created_at, updated_at, chicken, price_per_chicken, user_id, farmer_id)
VALUES($1, $2, $3, $4, $5, $6, $7)
RETURNING id, created_at, updated_at, chicken, price_per_chicken, user_id, farmer_id
`

type CreatePurchaseParams struct {
	ID              uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Chicken         int32
	PricePerChicken int32
	UserID          uuid.UUID
	FarmerID        uuid.UUID
}

func (q *Queries) CreatePurchase(ctx context.Context, arg CreatePurchaseParams) (Purchase, error) {
	row := q.db.QueryRowContext(ctx, createPurchase,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Chicken,
		arg.PricePerChicken,
		arg.UserID,
		arg.FarmerID,
	)
	var i Purchase
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Chicken,
		&i.PricePerChicken,
		&i.UserID,
		&i.FarmerID,
	)
	return i, err
}

const deletePurchase = `-- name: DeletePurchase :exec

DELETE FROM purchases WHERE id = $1
`

func (q *Queries) DeletePurchase(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deletePurchase, id)
	return err
}

const getMostRecentPurchase = `-- name: GetMostRecentPurchase :one
SELECT id, created_at, updated_at, chicken, price_per_chicken, user_id, farmer_id FROM purchases
ORDER BY created_at DESC
LIMIT 1
`

func (q *Queries) GetMostRecentPurchase(ctx context.Context) (Purchase, error) {
	row := q.db.QueryRowContext(ctx, getMostRecentPurchase)
	var i Purchase
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Chicken,
		&i.PricePerChicken,
		&i.UserID,
		&i.FarmerID,
	)
	return i, err
}

const getPurchaseByID = `-- name: GetPurchaseByID :one

SELECT id, created_at, updated_at, chicken, price_per_chicken, user_id, farmer_id FROM purchases
WHERE id = $1
`

func (q *Queries) GetPurchaseByID(ctx context.Context, id uuid.UUID) (Purchase, error) {
	row := q.db.QueryRowContext(ctx, getPurchaseByID, id)
	var i Purchase
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Chicken,
		&i.PricePerChicken,
		&i.UserID,
		&i.FarmerID,
	)
	return i, err
}

const getPurchases = `-- name: GetPurchases :many

SELECT id, created_at, updated_at, chicken, price_per_chicken, user_id, farmer_id FROM purchases ORDER BY created_at DESC
`

func (q *Queries) GetPurchases(ctx context.Context) ([]Purchase, error) {
	rows, err := q.db.QueryContext(ctx, getPurchases)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Purchase
	for rows.Next() {
		var i Purchase
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Chicken,
			&i.PricePerChicken,
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
