package api

import (
	"github.com/google/uuid"
	"github.com/mike-kimani/whitepointinventory/internal/purchases"
	"time"
)

type PurchaseResponse struct {
	ID                   uuid.UUID `json:"id"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	Chicken              int32     `json:"chicken"`
	PricePerChicken      int32     `json:"price_per_chicken"`
	UserID               uuid.UUID `json:"user_id"`
	FarmerID             uuid.UUID `json:"farmer_id"`
	UserName             string    `json:"user_name"`
	FarmerName           string    `json:"farmer_name"`
	FarmerChickenBalance float64   `json:"chicken_balance"`
	FarmerCashBalance    int32     `json:"cash_balance"`
}

func DomainPurchaseToPurchase(domainPurchase purchases.Purchase) PurchaseResponse {
	return PurchaseResponse{
		ID:                   domainPurchase.ID,
		CreatedAt:            domainPurchase.CreatedAt,
		UpdatedAt:            domainPurchase.UpdatedAt,
		Chicken:              domainPurchase.Chicken,
		PricePerChicken:      domainPurchase.PricePerChicken,
		UserID:               domainPurchase.UserID,
		FarmerID:             domainPurchase.FarmerID,
		UserName:             domainPurchase.UserName,
		FarmerName:           domainPurchase.FarmerName,
		FarmerChickenBalance: domainPurchase.FarmerChickenBalance,
		FarmerCashBalance:    domainPurchase.FarmerCashBalance,
	}
}

func DomainPurchasesToPurchases(domainPurchases []purchases.Purchase) []PurchaseResponse {
	var purchaseResponses []PurchaseResponse
	for _, domainPurchase := range domainPurchases {
		purchaseResponses = append(purchaseResponses, DomainPurchaseToPurchase(domainPurchase))
	}
	return purchaseResponses
}
