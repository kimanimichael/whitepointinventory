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
	user, err := s.repo.GetUserByID(ID)
	if err != nil {
		return nil, err
	}
	return &domain.User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Email:     user.Email,
		APIKey:    user.APIKey,
	}, nil
}

func (s *userService) GetUserByEmail(email string) (*domain.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return &domain.User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Email:     user.Email,
		APIKey:    user.APIKey,
		Password:  user.Password,
	}, nil
}

func (s *userService) GetUserByAPIKey(key string) (*domain.User, error) {
	user, err := s.repo.GetUserByAPIKey(key)
	if err != nil {
		return nil, err
	}
	return &domain.User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Email:     user.Email,
		APIKey:    user.APIKey,
		Password:  user.Password,
	}, nil
}

func (s *userService) GetUsers() ([]domain.User, error) {
	users, err := s.repo.GetUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}
