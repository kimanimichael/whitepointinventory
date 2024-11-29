package paymentsapi

type CreatePaymentRequest struct {
	CashPaid     int32  `json:"cash_paid"`
	ChickenPrice int32  `json:"price_per_chicken_paid"`
	FarmerName   string `json:"farmer_name"`
}

type GetTransactionRequest struct {
	ID string `json:"payment_id"`
}

type GetPagedPaymentsRequest struct {
	Offset uint32 `json:"offset"`
	Limit  uint32 `json:"limit"`
}
