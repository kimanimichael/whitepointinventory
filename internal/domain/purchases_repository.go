package domain

import "github.com/google/uuid"

type PurchasesRepository interface {
	CreatePurchase(chickenNo int32, chickenPrice int32, farmerName string, user *User) (*Purchase, error)
	GetPurchaseByID(ID uuid.UUID) (*Purchase, error)
	GetMostRecentPurchase() (*Purchase, error)
	GetPurchases() ([]Purchase, error)
	DeletePurchase(ID uuid.UUID) error
}
