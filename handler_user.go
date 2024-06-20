package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mike-kimani/whitepointinventory/internal/database"
)

const SecretKey = "secret"

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		Email    string `json:"email_address"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing json: %v", err))
		return
	}

	if params.Name == "" {
		respondWithError(w, 400, "Name field is empty")
		return
	}

	if params.Email == "" {
		respondWithError(w, 400, "Email field is empty")
		return
	}

	if len(params.Password) < 5 {
		respondWithError(w, 400, "Password field is empty or too short. Must be 5 or more characters")
		return
	}

	if !strings.Contains(params.Email, "@") {
		respondWithError(w, 400, "Invalid email format")
		return
	}

	if !strings.Contains(params.Email, ".") {
		respondWithError(w, 400, "Invalid email format")
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
		Password:  params.Password,
		Email:     params.Email,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}
	respondWithJSON(w, 201, user)
}