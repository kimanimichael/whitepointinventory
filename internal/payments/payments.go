package payments

import (
	"github.com/google/uuid"
	"github.com/mike-kimani/whitepointinventory/internal/users"
	"time"
)

type Payment struct {
	ID                   uuid.UUID
	CreatedAt            time.Time
	UpdatedAt            time.Time
	CashPaid             int32
	PricePerChickenPaid  int32
	UserID               uuid.UUID
	FarmerID             uuid.UUID
	UserName             string
	FarmerName           string
	FarmerChickenBalance float64
	FarmerCashBalance    int32
}

type PaymentsService interface {
	CreatePayment(cashPaid, chickenPrice int32, farmerName string, user *users.User) (*Payment, error)
	GetPaymentByID(ID uuid.UUID) (*Payment, error)
	GetPayments() ([]Payment, error)
	DeletePaymentByID(ID uuid.UUID) error
}

type PaymentsRepository interface {
	CreatePayment(cashPaid, chickenPrice int32, farmerName string, user *users.User) (*Payment, error)
	GetPaymentByID(ID uuid.UUID) (*Payment, error)
	GetMostRecentPayment() (*Payment, error)
	GetPayments() ([]Payment, error)
	DeletePayment(ID uuid.UUID) error
}
