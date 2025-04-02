package purchases

import (
	"context"
	"github.com/mike-kimani/whitepointinventory/internal/users"
	"time"
)

type Purchase struct {
	ID                   string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Chicken              int32
	PricePerChicken      int32
	UserID               string
	FarmerID             string
	UserName             string
	FarmerName           string
	FarmerChickenBalance float64
	FarmerCashBalance    int32
}

type page struct {
	Offset uint32
	Total  uint32
}

type PurchasePage struct {
	page
	Purchases []Purchase
}

type PurchaseRepository interface {
	CreatePurchase(ctx context.Context, chickenNo int32, chickenPrice int32, farmerName string, user *users.User) (*Purchase, error)
	GetPurchaseByID(ctx context.Context, ID string) (*Purchase, error)
	GetMostRecentPurchase(ctx context.Context) (*Purchase, error)
	GetPurchases(ctx context.Context) ([]Purchase, error)
	GetPagedPurchases(ctx context.Context, offset, limit uint32) (*PurchasePage, error)
	ChangePurchaseDate(ctx context.Context, purchaseID string, date time.Time, user *users.User) error
	DeletePurchase(ctx context.Context, ID string) error
}

type PurchaseService interface {
	CreatePurchase(ctx context.Context, chickenNo int32, chickenPrice int32, farmerName string, user *users.User) (*Purchase, error)
	GetPurchaseByID(ctx context.Context, ID string) (*Purchase, error)
	GetPurchases(ctx context.Context) ([]Purchase, error)
	GetPagedPurchases(ctx context.Context, offset, limit uint32) (*PurchasePage, error)
	ChangePurchaseDate(ctx context.Context, purchaseID string, date time.Time, user *users.User) error
	DeletePurchaseByID(ctx context.Context, ID string) error
}
