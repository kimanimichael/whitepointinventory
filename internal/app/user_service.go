package app

import (
	"errors"
	"github.com/google/uuid"
	"github.com/mike-kimani/whitepointinventory/internal/domain"
	"strings"
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
	if name == "" || email == "" || password == "" {
		return nil, errors.New("missing name, email or password")
	}
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return nil, errors.New("email must contain @ and . character")
	}
	user, err := s.repo.CreateUser(name, email, password)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		Name:  user.Name,
		Email: user.Email,
	}, nil
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
