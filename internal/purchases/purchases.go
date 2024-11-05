package purchases

import (
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
	CreatePurchase(chickenNo int32, chickenPrice int32, farmerName string, user *users.User) (*Purchase, error)
	GetPurchaseByID(ID uuid.UUID) (*Purchase, error)
	GetMostRecentPurchase() (*Purchase, error)
	GetPurchases() ([]Purchase, error)
	DeletePurchase(ID uuid.UUID) error
}

type PurchaseService interface {
	CreatePurchase(chickenNo int32, chickenPrice int32, farmerName string, user *users.User) (*Purchase, error)
	GetPurchaseByID(ID uuid.UUID) (*Purchase, error)
	GetPurchases() ([]Purchase, error)
	DeletePurchaseByID(ID uuid.UUID) error
}
