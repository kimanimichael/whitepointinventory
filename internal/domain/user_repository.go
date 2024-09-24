package domain

import "github.com/google/uuid"

type UserRepository interface {
	CreateUser(name, email, password string) (*User, error)
	GetUserByID(ID uuid.UUID) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUsers() ([]User, error)
}
