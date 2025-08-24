package usersapi

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mike-kimani/fechronizo/v2/pkg/httpresponses"
	"github.com/mike-kimani/whitepointinventory/internal/users"
	httpauth "github.com/mike-kimani/whitepointinventory/pkg/http"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

const secretKey = "secret"

type UserHandler struct {
	service users.UserService
}

func NewUserHandler(service users.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) RegisterRoutes(router chi.Router) {
	router.Post("/users", h.CreateUser)
	router.Post("/login", h.UserLogin)
	router.Get("/users", h.GetUserFromCookie)
	router.Get("/user", h.GetUsers)
	router.Post("/logout", h.UserLogout)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	params := CreateUserRequest{}

	decode := json.NewDecoder(r.Body)
	if err := decode.Decode(&params); err != nil {
		httpresponses.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), 14)
	if err != nil {
		httpresponses.RespondWithError(w, 500, fmt.Sprintf("Error hashing password: %v", err))
		return
	}
	hashedPasswordString := string(hashedPassword)

	ctx := r.Context()

	user, err := h.service.CreateUser(ctx, params.Name, params.Email, hashedPasswordString)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponses.RespondWithJson(w, http.StatusCreated, user)
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fetchedUsers, err := h.service.GetUsers(ctx)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponses.RespondWithJson(w, http.StatusOK, fetchedUsers)
}

func (h *UserHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	email, password, err := httpauth.GetPasswordAndEmailFromBody(r)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}

	ctx := r.Context()

	user, err := h.service.GetUserByEmail(ctx, email)
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

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    user.ID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
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
	userToReturn := userToResponseUser(*user)
	httpresponses.RespondWithJson(w, http.StatusOK, userToReturn)
}

func (h *UserHandler) GetUserFromCookie(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting cookie: %v", err))
		return
	}
	tokenString := cookie.Value
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting token: %v", err))
		return
	}
	claims := token.Claims.(*jwt.RegisteredClaims)
	userID := claims.Issuer

	ctx := r.Context()

	user, err := h.service.GetUserByID(ctx, userID)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting user: %v", err))
		return
	}
	userToReturn := userToResponseUser(*user)
	httpresponses.RespondWithJson(w, http.StatusOK, userToReturn)
}

func (h *UserHandler) UserLogout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(time.Hour * -1),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}
