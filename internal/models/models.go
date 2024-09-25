package models

import (
	"github.com/google/uuid"
	"github.com/mike-kimani/whitepointinventory/internal/adapters/database/sqlc/gensql"
	"time"
)

type ApiConfig struct {
	DB *sqlcdatabase.Queries
}

type Purchase struct {
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

type Payment struct {
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

type Farmer struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Name           string    `json:"name"`
	ChickenBalance float64   `json:"chicken_balance"`
	CashBalance    int32     `json:"cash_balance"`
}

type User struct {
	ID        uuid.UUID `json:"ID"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	Name      string    `json:"Name"`
	ApiKey    string    `json:"ApiKey"`
	Email     string    `json:"Email"`
}

func DatabasePurchaseToPurchase(dbPurchase sqlcdatabase.Purchase) Purchase {
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

func DatabasePurchasesToPurchases(dbPurchases []sqlcdatabase.Purchase) []Purchase {
	purchases := []Purchase{}
	for _, dbPurchase := range dbPurchases {
		purchases = append(purchases, DatabasePurchaseToPurchase(dbPurchase))
	}
	return purchases
}

func DatabasePaymentToPayment(dbPayment sqlcdatabase.Payment) Payment {
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

func DatabasePaymentsToPayments(dbPayments []sqlcdatabase.Payment) []Payment {
	payments := []Payment{}
	for _, dbPayment := range dbPayments {
		payments = append(payments, DatabasePaymentToPayment(dbPayment))
	}
	return payments
}

func DatabaseFarmerToFarmer(dbFarmer sqlcdatabase.Farmer) Farmer {
	return Farmer{
		ID:             dbFarmer.ID,
		CreatedAt:      dbFarmer.CreatedAt,
		UpdatedAt:      dbFarmer.UpdatedAt,
		Name:           dbFarmer.Name,
		ChickenBalance: dbFarmer.ChickenBalance.Float64,
		CashBalance:    dbFarmer.CashBalance.Int32,
	}
}

func DatabaseFarmersToFarmers(dbFarmers []sqlcdatabase.Farmer) []Farmer {
	farmers := []Farmer{}
	for _, dbFarmer := range dbFarmers {
		farmers = append(farmers, DatabaseFarmerToFarmer(dbFarmer))
	}
	return farmers
}

func DatabaseUserToUser(dbUser sqlcdatabase.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		ApiKey:    dbUser.ApiKey,
		Email:     dbUser.Email,
	}
}

func DatabaseUsersToUsers(dbUsers []sqlcdatabase.User) []User {
	var users []User
	for _, dbUser := range dbUsers {
		users = append(users, DatabaseUserToUser(dbUser))
	}
	return users
}
