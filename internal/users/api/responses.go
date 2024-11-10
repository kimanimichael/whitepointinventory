package usersapi

import (
	"github.com/google/uuid"
	"github.com/mike-kimani/whitepointinventory/internal/users"
	"time"
)

type userResponse struct {
	ID        uuid.UUID `json:"ID"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	Name      string    `json:"Name"`
	ApiKey    string    `json:"ApiKey"`
	Email     string    `json:"Email"`
	Password  string
}

func userToResponseUser(user users.User) userResponse {
	return userResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		ApiKey:    user.APIKey,
		Email:     user.Email,
	}
}

func usersToResponseUsers(users []users.User) []userResponse {
	var userResponses []userResponse
	for _, domainUser := range users {
		userResponses = append(userResponses, userToResponseUser(domainUser))
	}
	return userResponses
}
