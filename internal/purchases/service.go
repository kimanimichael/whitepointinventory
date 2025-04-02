package purchases

import (
	"context"
	"fmt"
	"github.com/mike-kimani/whitepointinventory/internal/users"
	"time"
)

// MinChickenNumber prevents zero value chicken number entries
const MinChickenNumber = 1

// MaxChickenNumber prevents entry errors e.g. 1000 instead of 100
const MaxChickenNumber = 999

// MinChickenPrice prevents entry errors e.g. swapping of chicken number and price
const MinChickenPrice = 100

// MaxChickenPrice prevents entry errors e.g. 4500 instead of 450 for price
const MaxChickenPrice = 1000

const IdenticalTransactionInterval = 2 * time.Minute

type purchaseService struct {
	repo PurchaseRepository
}

func NewPurchaseService(repo PurchaseRepository) PurchaseService {
	return &purchaseService{
		repo: repo,
	}
}

func (s *purchaseService) CreatePurchase(ctx context.Context, chickenNo, chickenPrice int32, farmerName string, user *users.User) (*Purchase, error) {
	if chickenNo < MinChickenNumber || chickenNo > MaxChickenNumber {
		return nil, fmt.Errorf("chicken number must be within %d and %d", MinChickenNumber, MaxChickenNumber)
	}

	if chickenPrice < MinChickenPrice || chickenPrice > MaxChickenPrice {
		return nil, fmt.Errorf("chicken price must be within %d and %d", MinChickenPrice, MaxChickenPrice)
	}
	mostRecentPurchase, err := s.repo.GetMostRecentPurchase(ctx)
	if err != nil {
		return nil, err
	}
	currentTime := time.Now()
	correctedMostRecentPurchaseTime := time.Date(
		mostRecentPurchase.CreatedAt.Year(),
		mostRecentPurchase.CreatedAt.Month(),
		mostRecentPurchase.CreatedAt.Day(),
		mostRecentPurchase.CreatedAt.Hour(),
		mostRecentPurchase.CreatedAt.Minute(),
		mostRecentPurchase.CreatedAt.Second(),
		mostRecentPurchase.CreatedAt.Nanosecond(),
		time.FixedZone("EAT", 3*60*60),
	)
	durationSinceLastPurchase := currentTime.Sub(correctedMostRecentPurchaseTime)

	if durationSinceLastPurchase < IdenticalTransactionInterval {
		fmt.Printf("Duration Since Last Payment less than %d minutes\n", IdenticalTransactionInterval/time.Minute)
		if mostRecentPurchase.FarmerName == farmerName {
			if mostRecentPurchase.Chicken == chickenNo {
				if mostRecentPurchase.PricePerChicken == chickenPrice {
					return nil, fmt.Errorf("identical purchase made for Farmer %s. Wait for %d seconds", farmerName, int(IdenticalTransactionInterval.Seconds()-durationSinceLastPurchase.Seconds()))
				}
			}
		}
	}

	purchase, err := s.repo.CreatePurchase(ctx, chickenNo, chickenPrice, farmerName, user)
	if err != nil {
		return nil, err
	}
	return purchase, nil
}

func (s *purchaseService) GetPurchaseByID(ctx context.Context, ID string) (*Purchase, error) {
	purchase, err := s.repo.GetPurchaseByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	return purchase, nil
}

func (s *purchaseService) GetPurchases(ctx context.Context) ([]Purchase, error) {
	purchases, err := s.repo.GetPurchases(ctx)
	if err != nil {
		return nil, err
	}
	return purchases, nil
}

func (s *purchaseService) GetPagedPurchases(ctx context.Context, offset, limit uint32) (*PurchasePage, error) {
	purchasePage, err := s.repo.GetPagedPurchases(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return purchasePage, nil
}

func (s *purchaseService) ChangePurchaseDate(ctx context.Context, paymentID string, date time.Time, user *users.User) error {
	err := s.repo.ChangePurchaseDate(ctx, paymentID, date, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *purchaseService) DeletePurchaseByID(ctx context.Context, ID string) error {
	err := s.repo.DeletePurchase(ctx, ID)
	if err != nil {
		return err
	}
	return nil
}
