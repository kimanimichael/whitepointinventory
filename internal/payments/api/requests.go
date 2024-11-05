package api

import "github.com/google/uuid"

type CreatePaymentRequest struct {
	CashPaid     int32  `json:"cash_paid"`
	ChickenPrice int32  `json:"price_per_chicken_paid"`
	FarmerName   string `json:"farmer_name"`
}

type GetTransactionRequest struct {
	ID uuid.UUID `json:"payment_id"`
}
