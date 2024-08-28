package middleware

import (
	"fmt"
	"github.com/mike-kimani/whitepointinventory/internal/database"
	"github.com/mike-kimani/whitepointinventory/pkg/auth"
	"github.com/mike-kimani/whitepointinventory/pkg/jsonresponses"
	"net/http"
)

type ApiConfig struct {
	DB *database.Queries
}

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *ApiConfig) MiddlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		APIKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Auth error: %v", err))
			return
		}
		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), APIKey)

		if err != nil {
			jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}
		handler(w, r, user)
	}
}
