package users

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mike-kimani/whitepointinventory/internal/adapters/database/sqlc/gensql"
	"time"
)

type UserRepositorySql struct {
	DB *sqlcdatabase.Queries
}

var _ userRepository = (*UserRepositorySql)(nil)

func NewUserRepositorySQL(db *sqlcdatabase.Queries) *UserRepositorySql {
	return &UserRepositorySql{
		DB: db,
	}
}

func (r *UserRepositorySql) CreateUser(name, email, password string) (*User, error) {
	user, err := r.DB.CreateUser(context.Background(), sqlcdatabase.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Email:     email,
		Password:  password,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating user: %v", err)
	}
	return &User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Email:     user.Email,
		APIKey:    user.ApiKey,
	}, nil
}

func (r *UserRepositorySql) GetUserByID(ID uuid.UUID) (*User, error) {
	user, err := r.DB.GetUserByID(context.Background(), ID)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}
	return &User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Email:     user.Email,
		APIKey:    user.ApiKey,
	}, nil
}

func (r *UserRepositorySql) GetUserByEmail(email string) (*User, error) {
	user, err := r.DB.GetUserByEmail(context.Background(), email)
	if err != nil {
		return nil, fmt.Errorf("error getting user from email: %v", err)
	}
	return &User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Email:     user.Email,
		APIKey:    user.ApiKey,
		Password:  user.Password,
	}, nil
}

func (r *UserRepositorySql) GetUserByAPIKey(key string) (*User, error) {
	user, err := r.DB.GetUserByAPIKey(context.Background(), key)
	if err != nil {
		return nil, fmt.Errorf("error getting user from APIKey: %v", err)
	}
	return &User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Email:     user.Email,
		APIKey:    user.ApiKey,
		Password:  user.Password,
	}, nil
}

func (r *UserRepositorySql) GetUsers() ([]User, error) {
	users, err := r.DB.GetUsers(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting users: %v", err)
	}

	var userList []User
	for _, user := range users {
		userList = append(userList, User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Name:      user.Name,
			Email:     user.Email,
		})
	}
	return userList, nil
}
