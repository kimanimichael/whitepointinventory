package purchasesapi

import "github.com/google/uuid"

type CreatePurchaseRequest struct {
	ChickenNo    int32  `json:"chicken_no"`
	ChickenPrice int32  `json:"chicken_price"`
	FarmerName   string `json:"farmer_name"`
}

type GetTransactionRequest struct {
	ID uuid.UUID `json:"payment_id"`
}
