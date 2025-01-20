package users

import (
	"context"
	"time"
)

type User struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Email     string
	APIKey    string
	Password  string
}

type UserService interface {
	CreateUser(ctx context.Context, name, email, password string) (*User, error)
	GetUserByID(ctx context.Context, ID string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByAPIKey(ctx context.Context, APIKey string) (*User, error)
	GetUsers(ctx context.Context) ([]User, error)
}

type userRepository interface {
	CreateUser(ctx context.Context, name, email, password string) (*User, error)
	GetUserByID(ctx context.Context, ID string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByAPIKey(ctx context.Context, key string) (*User, error)
	GetUsers(ctx context.Context) ([]User, error)
}
