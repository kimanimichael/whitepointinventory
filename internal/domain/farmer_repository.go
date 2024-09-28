package domain

import "github.com/google/uuid"

type FarmerRepository interface {
	CreateFarmer(name string, chickenBalance int32, cashBalance int32) (*Farmer, error)
	GetFarmerByName(string) (*Farmer, error)
	GetFarmers() ([]Farmer, error)
	DeleteFarmerByID(ID uuid.UUID) error
}
