package middleware

import (
	"fmt"
	"github.com/mike-kimani/whitepointinventory/internal/app"
	"github.com/mike-kimani/whitepointinventory/internal/domain"
	"github.com/mike-kimani/whitepointinventory/pkg/auth"
	"github.com/mike-kimani/whitepointinventory/pkg/jsonresponses"
	"net/http"
)

type userAuth struct {
	service app.UserService
}

type authedHandler func(http.ResponseWriter, *http.Request, *domain.User)

func (a *userAuth) MiddlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		APIKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Auth error: %v", err))
			return
		}
		user, err := a.service.GetUserByAPIKey(APIKey)

		if err != nil {
			jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}
		handler(w, r, user)
	}
}
