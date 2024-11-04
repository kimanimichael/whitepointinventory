package users

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Email     string
	APIKey    string
	Password  string
}

type UserService interface {
	CreateUser(name, email, password string) (*User, error)
	GetUserByID(ID uuid.UUID) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByAPIKey(APIKey string) (*User, error)
	GetUsers() ([]User, error)
}

type userRepository interface {
	CreateUser(name, email, password string) (*User, error)
	GetUserByID(ID uuid.UUID) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByAPIKey(key string) (*User, error)
	GetUsers() ([]User, error)
}
