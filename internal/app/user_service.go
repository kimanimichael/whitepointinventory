package app

import (
	"github.com/google/uuid"
	"github.com/mike-kimani/whitepointinventory/internal/domain"
)

type userService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) CreateUser(name, email, password string) (*domain.User, error) {
	return &domain.User{}, nil
}

func (s *userService) GetUserByID(ID uuid.UUID) (*domain.User, error) {
	return &domain.User{}, nil
}

func (s *userService) GetUserByEmail(email string) (*domain.User, error) {
	return &domain.User{}, nil
}

func (s *userService) GetUsers() ([]domain.User, error) {
	return []domain.User{}, nil
}
