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
