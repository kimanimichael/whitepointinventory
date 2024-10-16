package httpapi

import (
	"github.com/google/uuid"
	"github.com/mike-kimani/whitepointinventory/internal/domain"
	"time"
)

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
	Password  string
}

func DomainPurchaseToPurchase(domainPurchase domain.Purchase) Purchase {
	return Purchase{
		ID:                   domainPurchase.ID,
		CreatedAt:            domainPurchase.CreatedAt,
		UpdatedAt:            domainPurchase.UpdatedAt,
		Chicken:              domainPurchase.Chicken,
		PricePerChicken:      domainPurchase.PricePerChicken,
		UserID:               domainPurchase.UserID,
		FarmerID:             domainPurchase.FarmerID,
		UserName:             domainPurchase.UserName,
		FarmerName:           domainPurchase.FarmerName,
		FarmerChickenBalance: domainPurchase.FarmerChickenBalance,
		FarmerCashBalance:    domainPurchase.FarmerCashBalance,
	}
}

func DomainPurchasesToPurchases(domainPurchases []domain.Purchase) []Purchase {
	var purchases []Purchase
	for _, domainPurchase := range domainPurchases {
		purchases = append(purchases, DomainPurchaseToPurchase(domainPurchase))
	}
	return purchases
}

func DomainPaymentToPayment(domainPayment domain.Payment) Payment {
	return Payment{
		ID:                   domainPayment.ID,
		CreatedAt:            domainPayment.CreatedAt,
		UpdatedAt:            domainPayment.UpdatedAt,
		CashPaid:             domainPayment.CashPaid,
		PricePerChickenPaid:  domainPayment.PricePerChickenPaid,
		UserID:               domainPayment.UserID,
		FarmerID:             domainPayment.FarmerID,
		UserName:             domainPayment.UserName,
		FarmerName:           domainPayment.FarmerName,
		FarmerChickenBalance: domainPayment.FarmerChickenBalance,
		FarmerCashBalance:    domainPayment.FarmerCashBalance,
	}
}

func DomainPaymentsToPayments(domainPayments []domain.Payment) []Payment {
	var payments []Payment
	for _, domainPayment := range domainPayments {
		payments = append(payments, DomainPaymentToPayment(domainPayment))
	}
	return payments
}

func DomainFarmerToFarmer(domainFarmer domain.Farmer) Farmer {
	return Farmer{
		ID:             domainFarmer.ID,
		CreatedAt:      domainFarmer.CreatedAt,
		UpdatedAt:      domainFarmer.UpdatedAt,
		Name:           domainFarmer.Name,
		ChickenBalance: domainFarmer.ChickenBalance,
		CashBalance:    domainFarmer.CashBalance,
	}
}

func DomainFarmersToFarmers(domainFarmers []domain.Farmer) []Farmer {
	var farmers []Farmer
	for _, domainFarmer := range domainFarmers {
		farmers = append(farmers, DomainFarmerToFarmer(domainFarmer))
	}
	return farmers
}

func DomainUserToUser(domainUser domain.User) User {
	return User{
		ID:        domainUser.ID,
		CreatedAt: domainUser.CreatedAt,
		UpdatedAt: domainUser.UpdatedAt,
		Name:      domainUser.Name,
		ApiKey:    domainUser.APIKey,
		Email:     domainUser.Email,
	}
}

func DomainUsersToUsers(domainUsers []domain.User) []User {
	var users []User
	for _, domainUser := range domainUsers {
		users = append(users, DomainUserToUser(domainUser))
	}
	return users
}
