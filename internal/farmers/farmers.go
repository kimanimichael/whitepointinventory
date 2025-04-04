package farmers

import (
	"context"
	"time"
)

type Farmer struct {
	ID             string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Name           string
	PhoneNumber    string
	ChickenBalance float64
	CashBalance    int32
}

type Page struct {
	Offset uint32
	Total  uint32
}

type FarmersPage struct {
	Page
	Farmers []Farmer
}

type FarmerService interface {
	CreateFarmer(ctx context.Context, name string, chickenBalance float64, cashBalance int32) (*Farmer, error)
	GetFarmerByName(ctx context.Context, name string) (*Farmer, error)
	GetFarmers(ctx context.Context) ([]Farmer, error)
	GetPagedFarmers(ctx context.Context, offset, limit uint32) (*FarmersPage, error)
	SetFarmerBalances(ctx context.Context, name string, chickenBalance float64, cashBalance int32) (*Farmer, error)
	DeleteFarmerByID(ctx context.Context, ID string) error
}

type FarmerRepository interface {
	CreateFarmer(ctx context.Context, name string, chickenBalance float64, cashBalance int32) (*Farmer, error)
	GetFarmerByName(ctx context.Context, name string) (*Farmer, error)
	GetFarmers(ctx context.Context) ([]Farmer, error)
	GetPagedFarmers(ctx context.Context, offset, limit uint32) (*FarmersPage, error)
	SetFarmerBalances(ctx context.Context, name string, chickenBalance float64, cashBalance int32) (*Farmer, error)
	DeleteFarmerByID(ctx context.Context, ID string) error
}
