package domain

import "github.com/google/uuid"

type PaymentsRepository interface {
	CreatePayment(cashPaid, chickenPrice int32, farmerName string, user *User) (*Payment, error)
	GetPaymentByID(ID uuid.UUID) (*Payment, error)
	GetMostRecentPayment() (*Payment, error)
	GetPayments() ([]Payment, error)
	DeletePayment(ID uuid.UUID) error
}
