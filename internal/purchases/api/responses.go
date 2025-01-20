package purchasesapi

import (
	"github.com/mike-kimani/whitepointinventory/internal/purchases"
	"time"
)

type PurchaseResponse struct {
	ID                   string    `json:"id"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	Chicken              int32     `json:"chicken"`
	PricePerChicken      int32     `json:"price_per_chicken"`
	UserID               string    `json:"user_id"`
	FarmerID             string    `json:"farmer_id"`
	UserName             string    `json:"user_name"`
	FarmerName           string    `json:"farmer_name"`
	FarmerChickenBalance float64   `json:"chicken_balance"`
	FarmerCashBalance    int32     `json:"cash_balance"`
}

func purchaseToPurchaseResponse(purchase purchases.Purchase) PurchaseResponse {
	return PurchaseResponse{
		ID:                   purchase.ID,
		CreatedAt:            purchase.CreatedAt,
		UpdatedAt:            purchase.UpdatedAt,
		Chicken:              purchase.Chicken,
		PricePerChicken:      purchase.PricePerChicken,
		UserID:               purchase.UserID,
		FarmerID:             purchase.FarmerID,
		UserName:             purchase.UserName,
		FarmerName:           purchase.FarmerName,
		FarmerChickenBalance: purchase.FarmerChickenBalance,
		FarmerCashBalance:    purchase.FarmerCashBalance,
	}
}

func purchaseToPurchaseResponses(purchases []purchases.Purchase) []PurchaseResponse {
	var purchaseResponses []PurchaseResponse
	for _, domainPurchase := range purchases {
		purchaseResponses = append(purchaseResponses, purchaseToPurchaseResponse(domainPurchase))
	}
	return purchaseResponses
}
