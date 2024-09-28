package app

import (
	"github.com/google/uuid"
	"github.com/mike-kimani/whitepointinventory/internal/domain"
)

type UserService interface {
	CreateUser(name, email, password string) (*domain.User, error)
	GetUserByID(ID uuid.UUID) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	GetUsers() ([]domain.User, error)
}

type FarmerService interface {
	CreateFarmer(name string, chickenBalance int32, cashBalance int32) (*domain.Farmer, error)
	GetFarmerByName(name string) (*domain.Farmer, error)
	GetFarmers() ([]domain.Farmer, error)
	DeleteFarmerByID(ID uuid.UUID) error
}
