package app

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mike-kimani/whitepointinventory/internal/domain"
)

// MinCashPaid prevents swapping errors of cash paid and chicken price
const MinCashPaid = 1000

// MaxCashPaid prevents entry errors e.g. 400000 entry instead of 40000
const MaxCashPaid = 100000

//// MinChickenPrice prevents entry errors e.g. swapping of chicken number and price
//const MinChickenPrice = 100
//
//// MaxChickenPrice prevents entry errors e.g. 4500 instead of 450 for price
//const MaxChickenPrice = 1000

type paymentsService struct {
	repo domain.PaymentsRepository
}

func NewPaymentsService(repo domain.PaymentsRepository) PaymentsService {
	return &paymentsService{
		repo: repo,
	}
}

func (s *paymentsService) CreatePayment(cashPaid, chickenPrice int32, farmerName string, user *domain.User) (*domain.Payment, error) {
	if cashPaid < MinCashPaid || cashPaid > MaxCashPaid {
		return nil, fmt.Errorf("cash paid must be between %d and %d", MinCashPaid, MaxCashPaid)
	}
	if chickenPrice < MinChickenPrice || chickenPrice > MaxChickenPrice {
		return nil, fmt.Errorf("chicken price must be within %d and %d", MinChickenPrice, MaxChickenPrice)
	}

	payment, err := s.repo.CreatePayment(cashPaid, chickenPrice, farmerName, user)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *paymentsService) GetPaymentByID(ID uuid.UUID) (*domain.Payment, error) {
	payment, err := s.repo.GetPaymentByID(ID)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (s *paymentsService) GetPayments() ([]domain.Payment, error) {
	payments, err := s.repo.GetPayments()
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (s *paymentsService) DeletePaymentByID(ID uuid.UUID) error {
	err := s.repo.DeletePayment(ID)
	if err != nil {
		return err
	}
	return nil
}
