package sqlc

import (
	"context"
	"github.com/google/uuid"
	"github.com/mike-kimani/whitepointinventory/internal/adapters/database/sqlc/gensql"
	"github.com/mike-kimani/whitepointinventory/internal/domain"
	"github.com/mike-kimani/whitepointinventory/internal/models"
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
	modelUser := models.DatabaseUserToUser(user)
	return &domain.User{
		ID:        modelUser.ID,
		CreatedAt: modelUser.CreatedAt,
		UpdatedAt: modelUser.UpdatedAt,
		Name:      modelUser.Name,
		Email:     modelUser.Email,
		APIKey:    modelUser.ApiKey,
	}, nil
}

func (r *UserRepositorySql) GetUserByID(ID uuid.UUID) (*domain.User, error) {
	user, err := r.DB.GetUserByID(context.Background(), ID)
	if err != nil {
		return nil, err
	}
	modelUser := models.DatabaseUserToUser(user)
	return &domain.User{
		ID:        modelUser.ID,
		CreatedAt: modelUser.CreatedAt,
		UpdatedAt: modelUser.UpdatedAt,
		Name:      modelUser.Name,
		Email:     modelUser.Email,
		APIKey:    modelUser.ApiKey,
	}, nil
}

func (r *UserRepositorySql) GetUserByEmail(email string) (*domain.User, error) {
	user, err := r.DB.GetUserByEmail(context.Background(), email)
	if err != nil {
		return nil, err
	}
	modelUser := models.DatabaseUserToUser(user)
	return &domain.User{
		ID:        modelUser.ID,
		CreatedAt: modelUser.CreatedAt,
		UpdatedAt: modelUser.UpdatedAt,
		Name:      modelUser.Name,
		Email:     modelUser.Email,
		APIKey:    modelUser.ApiKey,
		Password:  modelUser.Password,
	}, nil
}

func (r *UserRepositorySql) GetUsers() ([]domain.User, error) {
	users, err := r.DB.GetUsers(context.Background())
	if err != nil {
		return nil, err
	}

	var userList []domain.User
	for _, user := range users {
		modelUser := models.DatabaseUserToUser(user)
		userList = append(userList, domain.User{
			ID:        modelUser.ID,
			CreatedAt: modelUser.CreatedAt,
			UpdatedAt: modelUser.UpdatedAt,
			Name:      modelUser.Name,
			Email:     modelUser.Email,
		})
	}
	return userList, nil
}
