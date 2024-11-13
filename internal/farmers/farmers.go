package farmers

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Farmer struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Name           string
	PhoneNumber    string
	ChickenBalance float64
	CashBalance    int32
}

type FarmerService interface {
	CreateFarmer(ctx context.Context, name string, chickenBalance float64, cashBalance int32) (*Farmer, error)
	GetFarmerByName(name string) (*Farmer, error)
	GetFarmers() ([]Farmer, error)
	DeleteFarmerByID(ID uuid.UUID) error
}

type FarmerRepository interface {
	CreateFarmer(ctx context.Context, name string, chickenBalance float64, cashBalance int32) (*Farmer, error)
	GetFarmerByName(string) (*Farmer, error)
	GetFarmers() ([]Farmer, error)
	DeleteFarmerByID(ID uuid.UUID) error
}
