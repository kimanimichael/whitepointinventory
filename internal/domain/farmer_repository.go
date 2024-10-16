package domain

import "github.com/google/uuid"

type FarmerRepository interface {
	CreateFarmer(name string, chickenBalance float64, cashBalance int32) (*Farmer, error)
	GetFarmerByName(string) (*Farmer, error)
	GetFarmers() ([]Farmer, error)
	DeleteFarmerByID(ID uuid.UUID) error
}
