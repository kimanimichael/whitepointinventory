package purchasesapi

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
