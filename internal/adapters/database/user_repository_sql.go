package sqlc

import (
	"context"
	"github.com/google/uuid"
	sqlcdatabase "github.com/mike-kimani/whitepointinventory/internal/adapters/db/sqlc/gensql"
	"github.com/mike-kimani/whitepointinventory/internal/domain"
	"time"
)

type UserRepositorySql struct {
	DB *sqlcdatabase.Queries
}

var _ domain.UserRepository = (*UserRepositorySql)(nil)

func (r *UserRepositorySql) CreateUser(name, email, password string) (*domain.User, error) {
	user, err := r.DB.CreateUser(context.Background(), sqlcdatabase.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Email:     email,
		Password:  password,
	})
	if err != nil {
		return nil, err
	}
	return &domain.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

func (r *UserRepositorySql) GetUserByID(ID uuid.UUID) (*domain.User, error) {
	user, err := r.DB.GetUserByID(context.Background(), ID)
	if err != nil {
		return nil, err
	}
	return &domain.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

func (r *UserRepositorySql) GetUserByEmail(email string) (*domain.User, error) {
	user, err := r.DB.GetUserByEmail(context.Background(), email)
	if err != nil {
		return nil, err
	}
	return &domain.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

func (r *UserRepositorySql) GetUsers() ([]domain.User, error) {
	users, err := r.DB.GetUsers(context.Background())
	if err != nil {
		return nil, err
	}

	var userList []domain.User
	for _, user := range users {
		userList = append(userList, domain.User{
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
		})
	}
	return userList, nil
}
