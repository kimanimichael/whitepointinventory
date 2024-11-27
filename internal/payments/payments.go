package payments

import (
	"context"
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

type page struct {
	Offset uint32
	Total  uint32
}

type PaymentPage struct {
	page
	Payments []Payment
}

type PaymentsService interface {
	CreatePayment(ctx context.Context, cashPaid, chickenPrice int32, farmerName string, user *users.User) (*Payment, error)
	GetPaymentByID(ctx context.Context, ID uuid.UUID) (*Payment, error)
	GetPayments(ctx context.Context) ([]Payment, error)
	GetPagedPayments(ctx context.Context, offset, limit uint32) (*PaymentPage, error)
	DeletePaymentByID(ctx context.Context, ID uuid.UUID) error
}

type PaymentsRepository interface {
	CreatePayment(ctx context.Context, cashPaid, chickenPrice int32, farmerName string, user *users.User) (*Payment, error)
	GetPaymentByID(ctx context.Context, ID uuid.UUID) (*Payment, error)
	GetMostRecentPayment(ctx context.Context) (*Payment, error)
	GetPayments(ctx context.Context) ([]Payment, error)
	GetPagedPayments(ctx context.Context, offset, limit uint32) (*PaymentPage, error)
	DeletePayment(ctx context.Context, ID uuid.UUID) error
}
