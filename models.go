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

type Farmer struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Name           string    `json:"name"`
	ChickenBalance int32     `json:"chicken_balance"`
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

func databaseFarmerToFarmer(dbFarmer database.Farmer) Farmer {
	return Farmer{
		ID:             dbFarmer.ID,
		CreatedAt:      dbFarmer.CreatedAt,
		UpdatedAt:      dbFarmer.UpdatedAt,
		Name:           dbFarmer.Name,
		ChickenBalance: dbFarmer.ChickenBalance.Int32,
		CashBalance:    dbFarmer.CashBalance.Int32,
	}
}

func databaseFarmersToFarmers(dbFarmers []database.Farmer) []Farmer {
	farmers := []Farmer{}
	for _, dbFarmer := range dbFarmers {
		farmers = append(farmers, databaseFarmerToFarmer(dbFarmer))
	}
	return farmers
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		ApiKey:    dbUser.ApiKey,
		Email:     dbUser.Email,
	}
}

func databaseUsersToUsers(dbUsers []database.User) []User {
	var users []User
	for _, dbUser := range dbUsers {
		users = append(users, databaseUserToUser(dbUser))
	}
	return users
}
