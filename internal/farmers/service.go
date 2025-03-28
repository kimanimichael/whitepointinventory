package farmers

import (
	"context"
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

func (s *farmerService) GetFarmerByName(ctx context.Context, name string) (*Farmer, error) {
	farmer, err := s.repo.GetFarmerByName(ctx, name)
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

func (s *farmerService) GetFarmers(ctx context.Context) ([]Farmer, error) {
	farmers, err := s.repo.GetFarmers(ctx)
	if err != nil {
		return nil, err
	}

	return farmers, nil
}

func (s *farmerService) GetPagedFarmers(ctx context.Context, offset, limit uint32) (*FarmersPage, error) {
	farmersPage, err := s.repo.GetPagedFarmers(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return farmersPage, nil
}

func (s *farmerService) SetFarmerBalances(ctx context.Context, name string, chickenBalance float64, cashBalance int32) (*Farmer, error) {
	farmer, err := s.repo.SetFarmerBalances(ctx, name, chickenBalance, cashBalance)
	if err != nil {
		return nil, err
	}
	return farmer, nil
}

func (s *farmerService) DeleteFarmerByID(ctx context.Context, ID string) error {
	err := s.repo.DeleteFarmerByID(ctx, ID)
	if err != nil {
		return err
	}
	return nil
}
