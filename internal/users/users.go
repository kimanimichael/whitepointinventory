package users

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mike-kimani/whitepointinventory/internal/models"
	"github.com/mike-kimani/whitepointinventory/pkg/auth"
	"github.com/mike-kimani/whitepointinventory/pkg/jsonresponses"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mike-kimani/whitepointinventory/internal/database"
)

type ApiConfig struct {
	DB *database.Queries
}

const SecretKey = "secret"

func (apiCfg *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		Email    string `json:"email_address"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Error parsing jsonresponses: %v", err))
		return
	}

	if params.Name == "" {
		jsonresponses.RespondWithError(w, 400, "Name field is empty")
		return
	}

	if params.Email == "" {
		jsonresponses.RespondWithError(w, 400, "Email field is empty")
		return
	}

	if len(params.Password) < 5 {
		jsonresponses.RespondWithError(w, 400, "Password field is empty or too short. Must be 5 or more characters")
		return
	}

	if !strings.Contains(params.Email, "@") {
		jsonresponses.RespondWithError(w, 400, "Invalid email format")
		return
	}

	if !strings.Contains(params.Email, ".") {
		jsonresponses.RespondWithError(w, 400, "Invalid email format")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), 14)
	if err != nil {
		jsonresponses.RespondWithError(w, 500, fmt.Sprintf("Error hashing password: %v", err))
		return
	}

	hashedPasswordString := string(hashedPassword)
	fmt.Println(hashedPasswordString)

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
		Password:  hashedPasswordString,
		Email:     params.Email,
	})

	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}
	jsonresponses.RespondWithJSON(w, 201, models.DatabaseUserToUser(user))
}

func (apiCfg *ApiConfig) HandlerUserLogin(w http.ResponseWriter, r *http.Request) {

	email, password, err := auth.GetPasswordAndEmailFromBody(r)
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Auth error: %v", err))
		return
	}
	user, err := apiCfg.DB.GetUserByEmail(r.Context(), email)
	if err != nil {
		jsonresponses.RespondWithError(w, 404, fmt.Sprintf("User does not exist: %v", err))
		return
	}
	userPasswordBytes := []byte(user.Password)

	err = bcrypt.CompareHashAndPassword(userPasswordBytes, []byte(password))
	if err != nil {
		jsonresponses.RespondWithError(w, 401, "Wrong password")
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.ID.String(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't get token: %v", err))
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	jsonresponses.RespondWithJSON(w, 200, models.DatabaseUserToUser(user))
}

func (apiCfg *ApiConfig) HandlerGetUserFromCookie(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		jsonresponses.RespondWithError(w, 404, fmt.Sprintf("Couldn't get cookie: %v", err))
		return
	}
	tokenString := cookie.Value
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		jsonresponses.RespondWithError(w, 404, fmt.Sprintf("Couldn't parse token: %v", err))
	}
	claims := token.Claims.(*jwt.StandardClaims)
	userID, err := uuid.Parse(claims.Issuer)

	user, err := apiCfg.DB.GetUserByID(r.Context(), userID)
	if err != nil {
		jsonresponses.RespondWithError(w, 404, fmt.Sprintf("User not found: %v", err))
		return
	}
	jsonresponses.RespondWithJSON(w, 200, models.DatabaseUserToUser(user))
}

func (apiCfg *ApiConfig) HandlerUserLogout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func (apiCfg *ApiConfig) HandlerGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := apiCfg.DB.GetUsers(r.Context())
	if err != nil {
		jsonresponses.RespondWithError(w, 500, fmt.Sprintf("Couldn't get users: %v", err))
		return
	}
	jsonresponses.RespondWithJSON(w, 200, models.DatabaseUsersToUsers(users))
}
