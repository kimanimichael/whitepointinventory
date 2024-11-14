package purchases

import (
	"context"
	"github.com/google/uuid"
	"github.com/mike-kimani/whitepointinventory/internal/users"
	"time"
)

type Purchase struct {
	ID                   uuid.UUID
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Chicken              int32
	PricePerChicken      int32
	UserID               uuid.UUID
	FarmerID             uuid.UUID
	UserName             string
	FarmerName           string
	FarmerChickenBalance float64
	FarmerCashBalance    int32
}

type PurchaseRepository interface {
	CreatePurchase(ctx context.Context, chickenNo int32, chickenPrice int32, farmerName string, user *users.User) (*Purchase, error)
	GetPurchaseByID(ctx context.Context, ID uuid.UUID) (*Purchase, error)
	GetMostRecentPurchase(ctx context.Context) (*Purchase, error)
	GetPurchases(ctx context.Context) ([]Purchase, error)
	DeletePurchase(ctx context.Context, ID uuid.UUID) error
}

type PurchaseService interface {
	CreatePurchase(ctx context.Context, chickenNo int32, chickenPrice int32, farmerName string, user *users.User) (*Purchase, error)
	GetPurchaseByID(ctx context.Context, ID uuid.UUID) (*Purchase, error)
	GetPurchases(ctx context.Context) ([]Purchase, error)
	DeletePurchaseByID(ctx context.Context, ID uuid.UUID) error
}
