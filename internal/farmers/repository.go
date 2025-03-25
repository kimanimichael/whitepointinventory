package farmers

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/mike-kimani/whitepointinventory/internal/adapters/database/sqlc/gensql"
	"time"
)

type FarmerRepositorySQL struct {
	DB *sqlcdatabase.Queries
}

var _ FarmerRepository = (*FarmerRepositorySQL)(nil)

func NewFarmerRepositorySQL(db *sqlcdatabase.Queries) *FarmerRepositorySQL {
	return &FarmerRepositorySQL{
		DB: db,
	}
}

func (r *FarmerRepositorySQL) CreateFarmer(ctx context.Context, name string, chickenBalance float64, cashBalance int32) (*Farmer, error) {
	var _chickenBalance sql.NullFloat64
	_chickenBalance.Float64 = chickenBalance
	_chickenBalance.Valid = true

	var _cashBalance sql.NullInt32
	_cashBalance.Int32 = cashBalance
	_cashBalance.Valid = true

	farmer, err := r.DB.CreateFarmer(ctx, sqlcdatabase.CreateFarmerParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Name:           name,
		ChickenBalance: _chickenBalance,
		CashBalance:    _cashBalance,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating farmer: %v", err)
	}
	return &Farmer{
		ID:             farmer.ID.String(),
		CreatedAt:      farmer.CreatedAt,
		UpdatedAt:      farmer.UpdatedAt,
		Name:           farmer.Name,
		ChickenBalance: farmer.ChickenBalance.Float64,
		CashBalance:    farmer.CashBalance.Int32,
	}, nil
}

func (r *FarmerRepositorySQL) GetFarmerByName(ctx context.Context, name string) (*Farmer, error) {
	farmer, err := r.DB.GetFarmerByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("error getting farmer: %v", err)
	}
	return &Farmer{
		ID:             farmer.ID.String(),
		CreatedAt:      farmer.CreatedAt,
		UpdatedAt:      farmer.UpdatedAt,
		Name:           farmer.Name,
		ChickenBalance: farmer.ChickenBalance.Float64,
		CashBalance:    farmer.CashBalance.Int32,
	}, nil
}

func (r *FarmerRepositorySQL) GetFarmers(ctx context.Context) ([]Farmer, error) {
	farmers, err := r.DB.GetFarmers(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating farmers: %v", err)
	}
	var farmersToReturn []Farmer
	for _, farmer := range farmers {
		farmersToReturn = append(farmersToReturn, Farmer{
			ID:             farmer.ID.String(),
			CreatedAt:      farmer.CreatedAt,
			UpdatedAt:      farmer.UpdatedAt,
			Name:           farmer.Name,
			ChickenBalance: farmer.ChickenBalance.Float64,
			CashBalance:    farmer.CashBalance.Int32,
		})
	}
	return farmersToReturn, nil
}

func (r *FarmerRepositorySQL) GetPagedFarmers(ctx context.Context, offset, limit uint32) (*FarmersPage, error) {
	farmers, err := r.DB.GetPagedFarmers(ctx, sqlcdatabase.GetPagedFarmersParams{
		Offset: int32(offset),
		Limit:  int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating farmers: %v", err)
	}
	total, err := r.DB.GetFarmerCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting farmer count: %v", err)
	}
	var farmersToReturn []Farmer
	for _, farmer := range farmers {
		farmersToReturn = append(farmersToReturn, Farmer{
			ID:             farmer.ID.String(),
			CreatedAt:      farmer.CreatedAt,
			UpdatedAt:      farmer.UpdatedAt,
			Name:           farmer.Name,
			ChickenBalance: farmer.ChickenBalance.Float64,
			CashBalance:    farmer.CashBalance.Int32,
		})
	}
	page := &FarmersPage{
		Page: Page{Offset: offset,
			Total: uint32(total),
		},
		Farmers: farmersToReturn,
	}
	return page, nil
}

func (r *FarmerRepositorySQL) SetFarmerBalances(ctx context.Context, name string, chickenBalance float64, cashBalance int32) (*Farmer, error) {
	var _chickenBalance sql.NullFloat64
	_chickenBalance.Float64 = chickenBalance
	_chickenBalance.Valid = true

	var _cashBalance sql.NullInt32
	_cashBalance.Int32 = cashBalance
	_cashBalance.Valid = true

	err := r.DB.SetFarmerBalances(ctx, sqlcdatabase.SetFarmerBalancesParams{
		Name:           name,
		ChickenBalance: _chickenBalance,
		CashBalance:    _cashBalance,
	})
	if err != nil {
		return nil, fmt.Errorf("error setting farmer balances: %v", err)
	}

	farmer, err := r.DB.GetFarmerByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("error getting farmer: %v", err)
	}
	return &Farmer{
		ID:             farmer.ID.String(),
		Name:           farmer.Name,
		ChickenBalance: farmer.ChickenBalance.Float64,
		CashBalance:    farmer.CashBalance.Int32,
	}, nil
}

func (r *FarmerRepositorySQL) DeleteFarmerByID(ctx context.Context, ID string) error {
	userID, err := uuid.Parse(ID)
	if err != nil {
		return fmt.Errorf("error parsing farmer ID: %v", err)
	}
	err = r.DB.DeleteFarmers(ctx, userID)
	if err != nil {
		return fmt.Errorf("error deleting farmer: %v", err)
	}
	return nil
}
