package app

import (
	"github.com/google/uuid"
	"github.com/mike-kimani/whitepointinventory/internal/domain"
)

type farmerService struct {
	repo domain.FarmerRepository
}

func NewFarmerService(repo domain.FarmerRepository) FarmerService {
	return &farmerService{
		repo: repo,
	}
}

func (s *farmerService) CreateFarmer(name string, chickenBalance float64, cashBalance int32) (*domain.Farmer, error) {
	farmer, err := s.repo.CreateFarmer(name, chickenBalance, cashBalance)
	if err != nil {
		return nil, err
	}
	return &domain.Farmer{
		ID:             farmer.ID,
		CreatedAt:      farmer.CreatedAt,
		UpdatedAt:      farmer.UpdatedAt,
		Name:           farmer.Name,
		ChickenBalance: farmer.ChickenBalance,
		CashBalance:    farmer.CashBalance,
	}, nil
}

func (s *farmerService) GetFarmerByName(name string) (*domain.Farmer, error) {
	farmer, err := s.repo.GetFarmerByName(name)
	if err != nil {
		return nil, err
	}
	return &domain.Farmer{
		ID:             farmer.ID,
		CreatedAt:      farmer.CreatedAt,
		UpdatedAt:      farmer.UpdatedAt,
		Name:           farmer.Name,
		ChickenBalance: farmer.ChickenBalance,
		CashBalance:    farmer.CashBalance,
	}, nil
}

func (s *farmerService) GetFarmers() ([]domain.Farmer, error) {
	farmers, err := s.repo.GetFarmers()
	if err != nil {
		return nil, err
	}

	return farmers, nil
}

func (s *farmerService) DeleteFarmerByID(ID uuid.UUID) error {
	err := s.repo.DeleteFarmerByID(ID)
	if err != nil {
		return err
	}
	return nil
}
