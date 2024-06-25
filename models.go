package main

import (
	"github.com/google/uuid"
	"github.com/mike-kimani/whitepointinventory/internal/database"
	"time"
)

type Purchase struct {
	ID              uuid.UUID `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Chicken         int32     `json:"chicken"`
	PricePerChicken int32     `json:"price_per_chicken"`
	UserID          uuid.UUID `json:"user_id"`
	FarmerID        uuid.UUID `json:"farmer_id"`
	UserName        string    `json:"user_name"`
	FarmerName      string    `json:"farmer_name"`
}

type Payment struct {
	ID                  uuid.UUID `json:"id"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	CashPaid            int32     `json:"cash_paid"`
	PricePerChickenPaid int32     `json:"price_per_chicken_paid"`
	UserID              uuid.UUID `json:"user_id"`
	FarmerID            uuid.UUID `json:"farmer_id"`
	UserName            string    `json:"user_name"`
	FarmerName          string    `json:"farmer_name"`
}

func databasePurchaseToPurchase(dbPurchase database.Purchase) Purchase {
	return Purchase{
		ID:              dbPurchase.ID,
		CreatedAt:       dbPurchase.CreatedAt,
		UpdatedAt:       dbPurchase.UpdatedAt,
		Chicken:         dbPurchase.Chicken,
		PricePerChicken: dbPurchase.PricePerChicken,
		UserID:          dbPurchase.UserID,
		FarmerID:        dbPurchase.FarmerID,
	}
}

func databasePurchasesToPurchases(dbPurchases []database.Purchase) []Purchase {
	purchases := []Purchase{}
	for _, dbPurchase := range dbPurchases {
		purchases = append(purchases, databasePurchaseToPurchase(dbPurchase))
	}
	return purchases
}

func databasePaymentToPayment(dbPayment database.Payment) Payment {
	return Payment{
		ID:                  dbPayment.ID,
		CreatedAt:           dbPayment.CreatedAt,
		UpdatedAt:           dbPayment.UpdatedAt,
		CashPaid:            dbPayment.CashPaid,
		PricePerChickenPaid: dbPayment.PricePerChickenPaid,
		UserID:              dbPayment.UserID,
		FarmerID:            dbPayment.FarmerID,
	}
}

func databasePaymentsToPayments(dbPayments []database.Payment) []Payment {
	payments := []Payment{}
	for _, dbPayment := range dbPayments {
		payments = append(payments, databasePaymentToPayment(dbPayment))
	}
	return payments
}
