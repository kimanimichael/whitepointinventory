// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: farmers.sql

package sqlcdatabase

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createFarmer = `-- name: CreateFarmer :one
INSERT INTO farmers(id, created_at, updated_at, name, chicken_balance, cash_balance)
VALUES($1, $2, $3, $4, $5, $6)
RETURNING id, created_at, updated_at, name, chicken_balance, cash_balance
`

type CreateFarmerParams struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Name           string
	ChickenBalance sql.NullFloat64
	CashBalance    sql.NullInt32
}

func (q *Queries) CreateFarmer(ctx context.Context, arg CreateFarmerParams) (Farmer, error) {
	row := q.db.QueryRowContext(ctx, createFarmer,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
		arg.ChickenBalance,
		arg.CashBalance,
	)
	var i Farmer
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.ChickenBalance,
		&i.CashBalance,
	)
	return i, err
}

const decreaseCashOwed = `-- name: DecreaseCashOwed :exec
UPDATE farmers
SET cash_balance = COALESCE(cash_balance, 0) - ($1)
WHERE farmers.id = $2
`

type DecreaseCashOwedParams struct {
	CashBalance sql.NullInt32
	ID          uuid.UUID
}

func (q *Queries) DecreaseCashOwed(ctx context.Context, arg DecreaseCashOwedParams) error {
	_, err := q.db.ExecContext(ctx, decreaseCashOwed, arg.CashBalance, arg.ID)
	return err
}

const decreaseChickenOwed = `-- name: DecreaseChickenOwed :exec
UPDATE farmers
SET chicken_balance = COALESCE(chicken_balance, 0) - $1
WHERE farmers.id = $2
`

type DecreaseChickenOwedParams struct {
	ChickenBalance sql.NullFloat64
	ID             uuid.UUID
}

func (q *Queries) DecreaseChickenOwed(ctx context.Context, arg DecreaseChickenOwedParams) error {
	_, err := q.db.ExecContext(ctx, decreaseChickenOwed, arg.ChickenBalance, arg.ID)
	return err
}

const deleteFarmers = `-- name: DeleteFarmers :exec
DELETE FROM farmers where id = $1
`

func (q *Queries) DeleteFarmers(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteFarmers, id)
	return err
}

const getFarmerByID = `-- name: GetFarmerByID :one
SELECT id, created_at, updated_at, name, chicken_balance, cash_balance FROM farmers where id = $1
`

func (q *Queries) GetFarmerByID(ctx context.Context, id uuid.UUID) (Farmer, error) {
	row := q.db.QueryRowContext(ctx, getFarmerByID, id)
	var i Farmer
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.ChickenBalance,
		&i.CashBalance,
	)
	return i, err
}

const getFarmerByName = `-- name: GetFarmerByName :one
SELECT id, created_at, updated_at, name, chicken_balance, cash_balance FROM farmers where name = $1
`

func (q *Queries) GetFarmerByName(ctx context.Context, name string) (Farmer, error) {
	row := q.db.QueryRowContext(ctx, getFarmerByName, name)
	var i Farmer
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.ChickenBalance,
		&i.CashBalance,
	)
	return i, err
}

const getFarmerCount = `-- name: GetFarmerCount :one
SELECT COUNT(*) AS total FROM farmers
`

func (q *Queries) GetFarmerCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, getFarmerCount)
	var total int64
	err := row.Scan(&total)
	return total, err
}

const getFarmers = `-- name: GetFarmers :many
SELECT id, created_at, updated_at, name, chicken_balance, cash_balance FROM farmers ORDER BY updated_at DESC
`

func (q *Queries) GetFarmers(ctx context.Context) ([]Farmer, error) {
	rows, err := q.db.QueryContext(ctx, getFarmers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Farmer
	for rows.Next() {
		var i Farmer
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.ChickenBalance,
			&i.CashBalance,
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

const getPagedFarmers = `-- name: GetPagedFarmers :many
SELECT id, created_at, updated_at, name, chicken_balance, cash_balance FROM farmers OFFSET $1 LIMIT $2
`

type GetPagedFarmersParams struct {
	Offset int32
	Limit  int32
}

func (q *Queries) GetPagedFarmers(ctx context.Context, arg GetPagedFarmersParams) ([]Farmer, error) {
	rows, err := q.db.QueryContext(ctx, getPagedFarmers, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Farmer
	for rows.Next() {
		var i Farmer
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.ChickenBalance,
			&i.CashBalance,
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

const increaseCashOwed = `-- name: IncreaseCashOwed :exec
UPDATE farmers
SET cash_balance = COALESCE(cash_balance, 0) + ($1)
WHERE farmers.id = $2
`

type IncreaseCashOwedParams struct {
	CashBalance sql.NullInt32
	ID          uuid.UUID
}

func (q *Queries) IncreaseCashOwed(ctx context.Context, arg IncreaseCashOwedParams) error {
	_, err := q.db.ExecContext(ctx, increaseCashOwed, arg.CashBalance, arg.ID)
	return err
}

const increaseChickenOwed = `-- name: IncreaseChickenOwed :exec
UPDATE farmers
SET chicken_balance = COALESCE(chicken_balance, 0) + $1
WHERE farmers.id = $2
`

type IncreaseChickenOwedParams struct {
	ChickenBalance sql.NullFloat64
	ID             uuid.UUID
}

func (q *Queries) IncreaseChickenOwed(ctx context.Context, arg IncreaseChickenOwedParams) error {
	_, err := q.db.ExecContext(ctx, increaseChickenOwed, arg.ChickenBalance, arg.ID)
	return err
}

const markFarmerAsUpdated = `-- name: MarkFarmerAsUpdated :exec
UPDATE farmers
SET updated_at = NOW()
WHERE id = $1
`

func (q *Queries) MarkFarmerAsUpdated(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, markFarmerAsUpdated, id)
	return err
}

const setFarmerBalances = `-- name: SetFarmerBalances :exec
UPDATE farmers
SET chicken_balance = $2, cash_balance = $3
where farmers.name = $1
`

type SetFarmerBalancesParams struct {
	Name           string
	ChickenBalance sql.NullFloat64
	CashBalance    sql.NullInt32
}

func (q *Queries) SetFarmerBalances(ctx context.Context, arg SetFarmerBalancesParams) error {
	_, err := q.db.ExecContext(ctx, setFarmerBalances, arg.Name, arg.ChickenBalance, arg.CashBalance)
	return err
}
