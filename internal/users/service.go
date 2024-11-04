package users

import (
	"errors"
	"github.com/google/uuid"
	"strings"
)

type service struct {
	repo userRepository
}

func NewUserService(repo userRepository) UserService {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateUser(name, email, password string) (*User, error) {
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

	return &User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Email:     user.Email,
		APIKey:    user.APIKey,
	}, nil
}

func (s *service) GetUserByID(ID uuid.UUID) (*User, error) {
	user, err := s.repo.GetUserByID(ID)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Email:     user.Email,
		APIKey:    user.APIKey,
	}, nil
}

func (s *service) GetUserByEmail(email string) (*User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Email:     user.Email,
		APIKey:    user.APIKey,
		Password:  user.Password,
	}, nil
}

func (s *service) GetUserByAPIKey(key string) (*User, error) {
	user, err := s.repo.GetUserByAPIKey(key)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Email:     user.Email,
		APIKey:    user.APIKey,
		Password:  user.Password,
	}, nil
}

func (s *service) GetUsers() ([]User, error) {
	users, err := s.repo.GetUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}
