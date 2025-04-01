package payments

import (
	"context"
	"fmt"
	"github.com/mike-kimani/whitepointinventory/internal/users"
	"log"
	"time"
)

// MinCashPaid prevents swapping errors of cash paid and chicken price
const MinCashPaid = 1000

// MaxCashPaid prevents entry errors e.g. 400000 entry instead of 40000
const MaxCashPaid = 100000

// MinChickenPrice prevents entry errors e.g. swapping of chicken number and price
const MinChickenPrice = 100

// MaxChickenPrice prevents entry errors e.g. 4500 instead of 450 for price
const MaxChickenPrice = 1000

const IdenticalTransactionInterval = 2 * time.Minute

//// MinChickenPrice prevents entry errors e.g. swapping of chicken number and price
//const MinChickenPrice = 100
//
//// MaxChickenPrice prevents entry errors e.g. 4500 instead of 450 for price
//const MaxChickenPrice = 1000

type paymentsService struct {
	repo PaymentsRepository
}

func NewPaymentsService(repo PaymentsRepository) PaymentsService {
	return &paymentsService{
		repo: repo,
	}
}

func (s *paymentsService) CreatePayment(ctx context.Context, cashPaid, chickenPrice int32, farmerName string, user *users.User) (*Payment, error) {
	if cashPaid < MinCashPaid || cashPaid > MaxCashPaid {
		return nil, fmt.Errorf("cash paid must be between %d and %d", MinCashPaid, MaxCashPaid)
	}
	if chickenPrice < MinChickenPrice || chickenPrice > MaxChickenPrice {
		return nil, fmt.Errorf("chicken price must be within %d and %d", MinChickenPrice, MaxChickenPrice)
	}

	mostRecentPayment, err := s.repo.GetMostRecentPayment(ctx)
	if err != nil {
		log.Println("Could not get most recent payment:", err)
	} else {
		currentTime := time.Now()
		correctedMostRecentPaymentTime := time.Date(
			mostRecentPayment.CreatedAt.Year(),
			mostRecentPayment.CreatedAt.Month(),
			mostRecentPayment.CreatedAt.Day(),
			mostRecentPayment.CreatedAt.Hour(),
			mostRecentPayment.CreatedAt.Minute(),
			mostRecentPayment.CreatedAt.Second(),
			mostRecentPayment.CreatedAt.Nanosecond(),
			time.FixedZone("EAT", 3*60*60),
		)
		durationSinceLastPayment := currentTime.Sub(correctedMostRecentPaymentTime)

		if durationSinceLastPayment < IdenticalTransactionInterval {
			fmt.Printf("Duration Since Last Payment less than %d minutes\n", IdenticalTransactionInterval/time.Minute)
			if mostRecentPayment.FarmerName == farmerName {
				if mostRecentPayment.CashPaid == cashPaid {
					if mostRecentPayment.PricePerChickenPaid == chickenPrice {
						return nil, fmt.Errorf("identical payment made for Farmer %s. Wait for %d seconds", farmerName, int(IdenticalTransactionInterval.Seconds()-durationSinceLastPayment.Seconds()))
					}
				}
			}
		}
	}

	payment, err := s.repo.CreatePayment(ctx, cashPaid, chickenPrice, farmerName, user)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *paymentsService) GetPaymentByID(ctx context.Context, ID string) (*Payment, error) {
	payment, err := s.repo.GetPaymentByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (s *paymentsService) GetPayments(ctx context.Context) ([]Payment, error) {
	payments, err := s.repo.GetPayments(ctx)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (s *paymentsService) GetPagedPayments(ctx context.Context, offset, limit uint32) (*PaymentPage, error) {
	payments, err := s.repo.GetPagedPayments(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (s *paymentsService) DeletePaymentByID(ctx context.Context, ID string) error {
	err := s.repo.DeletePayment(ctx, ID)
	if err != nil {
		return err
	}
	return nil
}
