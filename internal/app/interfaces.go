package app

import (
	"github.com/google/uuid"
	"github.com/mike-kimani/whitepointinventory/internal/domain"
)

type UserService interface {
	CreateUser(name, email, password string) (*domain.User, error)
	GetUserByID(ID uuid.UUID) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	GetUserByAPIKey(APIKey string) (*domain.User, error)
	GetUsers() ([]domain.User, error)
}

type FarmerService interface {
	CreateFarmer(name string, chickenBalance float64, cashBalance int32) (*domain.Farmer, error)
	GetFarmerByName(name string) (*domain.Farmer, error)
	GetFarmers() ([]domain.Farmer, error)
	DeleteFarmerByID(ID uuid.UUID) error
}

type PurchaseService interface {
	CreatePurchase(chickenNo int32, chickenPrice int32, farmerName string, user *domain.User) (*domain.Purchase, error)
	GetPurchaseByID(ID uuid.UUID) (*domain.Purchase, error)
	GetPurchases() ([]domain.Purchase, error)
	DeletePurchaseByID(ID uuid.UUID) error
}

type PaymentsService interface {
	CreatePayment(cashPaid, chickenPrice int32, farmerName string, user *domain.User) (*domain.Payment, error)
	GetPaymentByID(ID uuid.UUID) (*domain.Payment, error)
	GetPayments() ([]domain.Payment, error)
	DeletePaymentByID(ID uuid.UUID) error
}
