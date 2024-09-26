package httpapi

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/mike-kimani/fechronizo/v2/pkg/httpresponses"
	"github.com/mike-kimani/whitepointinventory/internal/app"
	"github.com/mike-kimani/whitepointinventory/pkg/auth"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

const secretKey = "emaakama"

type UserHandler struct {
	service app.UserService
}

func NewUserHandler(service app.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) RegisterRoutes(router chi.Router) {
	router.Post("/user", h.CreateUser)
	router.Post("/login", h.UserLogin)
	router.Get("/users", h.GetUserFromCookie)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	params := CreateUserRequest{}

	decode := json.NewDecoder(r.Body)
	if err := decode.Decode(&params); err != nil {
		httpresponses.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	user, err := h.service.CreateUser(params.Name, params.Email, params.Password)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponses.RespondWithJson(w, http.StatusCreated, user)
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUsers()
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponses.RespondWithJson(w, http.StatusOK, users)
}

func (h *UserHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	email, password, err := auth.GetPasswordAndEmailFromBody(r)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}

	user, err := h.service.GetUserByEmail(email)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("User not found: %v", err))
		return
	}

	userPasswordBytes := []byte(user.Password)
	err = bcrypt.CompareHashAndPassword(userPasswordBytes, []byte(password))
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Wrong password: %v", err))
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.ID.String(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
	token, err := claims.SignedString([]byte(secretKey))
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting token: %v", err))
		return
	}
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	httpresponses.RespondWithJson(w, http.StatusOK, user)
}

func (h *UserHandler) GetUserFromCookie(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting cookie: %v", err))
	}
	tokenString := cookie.Value
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting token: %v", err))
	}
	claims := token.Claims.(*jwt.StandardClaims)
	userID, err := uuid.Parse(claims.Issuer)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error parsing userID: %v", err))
	}
	user, err := h.service.GetUserByID(userID)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting user: %v", err))
	}
	httpresponses.RespondWithJson(w, http.StatusOK, user)
}
