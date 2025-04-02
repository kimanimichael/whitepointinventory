package purchasesapi

import "time"

type CreatePurchaseRequest struct {
	ChickenNo    int32  `json:"chicken_no"`
	ChickenPrice int32  `json:"chicken_price"`
	FarmerName   string `json:"farmer_name"`
}

type GetTransactionRequest struct {
	ID string `json:"payment_id"`
}

type GetPagedPurchasesRequest struct {
	Offset uint32 `json:"offset"`
	Limit  uint32 `json:"limit"`
}

type ChangePurchaseDateRequest struct {
	ID   string    `json:"purchase_id"`
	Time time.Time `json:"new_time"`
}
