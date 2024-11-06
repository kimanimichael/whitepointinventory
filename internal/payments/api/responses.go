package paymentsapi

import (
	"github.com/google/uuid"
	"github.com/mike-kimani/whitepointinventory/internal/payments"
	"time"
)

type PaymentResponse struct {
	ID                   uuid.UUID `json:"id"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	CashPaid             int32     `json:"cash_paid"`
	PricePerChickenPaid  int32     `json:"price_per_chicken_paid"`
	UserID               uuid.UUID `json:"user_id"`
	FarmerID             uuid.UUID `json:"farmer_id"`
	UserName             string    `json:"user_name"`
	FarmerName           string    `json:"farmer_name"`
	FarmerChickenBalance float64   `json:"chicken_balance"`
	FarmerCashBalance    int32     `json:"cash_balance"`
}

func paymentToPaymentResponse(payment payments.Payment) PaymentResponse {
	return PaymentResponse{
		ID:                   payment.ID,
		CreatedAt:            payment.CreatedAt,
		UpdatedAt:            payment.UpdatedAt,
		CashPaid:             payment.CashPaid,
		PricePerChickenPaid:  payment.PricePerChickenPaid,
		UserID:               payment.UserID,
		FarmerID:             payment.FarmerID,
		UserName:             payment.UserName,
		FarmerName:           payment.FarmerName,
		FarmerChickenBalance: payment.FarmerChickenBalance,
		FarmerCashBalance:    payment.FarmerCashBalance,
	}
}

func paymentsToPaymentsResponses(domainPayments []payments.Payment) []PaymentResponse {
	var payments []PaymentResponse
	for _, domainPayment := range domainPayments {
		payments = append(payments, paymentToPaymentResponse(domainPayment))
	}
	return payments
}
