package farmers

import (
	"context"
	"github.com/google/uuid"
)

type farmerService struct {
	repo FarmerRepository
}

func NewFarmerService(repo FarmerRepository) FarmerService {
	return &farmerService{
		repo: repo,
	}
}

func (s *farmerService) CreateFarmer(ctx context.Context, name string, chickenBalance float64, cashBalance int32) (*Farmer, error) {
	farmer, err := s.repo.CreateFarmer(ctx, name, chickenBalance, cashBalance)
	if err != nil {
		return nil, err
	}
	return &Farmer{
		ID:             farmer.ID,
		CreatedAt:      farmer.CreatedAt,
		UpdatedAt:      farmer.UpdatedAt,
		Name:           farmer.Name,
		ChickenBalance: farmer.ChickenBalance,
		CashBalance:    farmer.CashBalance,
	}, nil
}

func (s *farmerService) GetFarmerByName(name string) (*Farmer, error) {
	farmer, err := s.repo.GetFarmerByName(name)
	if err != nil {
		return nil, err
	}
	return &Farmer{
		ID:             farmer.ID,
		CreatedAt:      farmer.CreatedAt,
		UpdatedAt:      farmer.UpdatedAt,
		Name:           farmer.Name,
		ChickenBalance: farmer.ChickenBalance,
		CashBalance:    farmer.CashBalance,
	}, nil
}

func (s *farmerService) GetFarmers() ([]Farmer, error) {
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
