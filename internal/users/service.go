package users

import (
	"context"
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

func (s *service) CreateUser(ctx context.Context, name, email, password string) (*User, error) {
	if name == "" || email == "" || password == "" {
		return nil, errors.New("missing name, email or password")
	}
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return nil, errors.New("email must contain @ and . character")
	}
	user, err := s.repo.CreateUser(ctx, name, email, password)
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

func (s *service) GetUserByID(ctx context.Context, ID uuid.UUID) (*User, error) {
	user, err := s.repo.GetUserByID(ctx, ID)
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

func (s *service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
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

func (s *service) GetUserByAPIKey(ctx context.Context, key string) (*User, error) {
	user, err := s.repo.GetUserByAPIKey(ctx, key)
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

func (s *service) GetUsers(ctx context.Context) ([]User, error) {
	users, err := s.repo.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}
