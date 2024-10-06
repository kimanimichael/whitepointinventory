package app

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mike-kimani/whitepointinventory/internal/domain"
)

// MinChickenNumber prevents zero value chicken number entries
const MinChickenNumber = 1

// MaxChickenNumber prevents entry errors e.g. 1000 instead of 100
const MaxChickenNumber = 999

// MinChickenPrice prevents entry errors e.g. swapping of chicken number and price
const MinChickenPrice = 100

// MaxChickenPrice prevents entry errors e.g. 4500 instead of 450 for price
const MaxChickenPrice = 1000

type purchaseService struct {
	repo domain.PurchasesRepository
}

func NewPurchaseService(repo domain.PurchasesRepository) PurchaseService {
	return &purchaseService{
		repo: repo,
	}
}

func (s *purchaseService) CreatePurchase(chickenNo, chickenPrice int32, farmerName string, user *domain.User) (*domain.Purchase, error) {
	if chickenNo < MinChickenNumber || chickenNo > MaxChickenNumber {
		return nil, fmt.Errorf("chicken number must be within %d and %d", MinChickenNumber, MaxChickenNumber)
	}

	if chickenPrice < MinChickenPrice || chickenPrice > MaxChickenPrice {
		return nil, fmt.Errorf("chicken price must be within %d and %d", MinChickenPrice, MaxChickenPrice)
	}
	purchase, err := s.repo.CreatePurchase(chickenNo, chickenPrice, farmerName, user)
	if err != nil {
		return nil, err
	}
	return purchase, nil
}

func (s *purchaseService) GetPurchaseByID(ID uuid.UUID) (*domain.Purchase, error) {
	purchase, err := s.repo.GetPurchaseByID(ID)
	if err != nil {
		return nil, err
	}
	return purchase, nil
}

func (s *purchaseService) GetPurchases() ([]domain.Purchase, error) {
	purchases, err := s.repo.GetPurchases()
	if err != nil {
		return nil, err
	}
	return purchases, nil
}

func (s *purchaseService) DeletePurchaseByID(ID uuid.UUID) error {
	err := s.repo.DeletePurchase(ID)
	if err != nil {
		return err
	}
	return nil
}
