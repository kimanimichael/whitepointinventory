package middleware

import (
	"fmt"
	"github.com/mike-kimani/fechronizo/v2/pkg/httpresponses"
	"github.com/mike-kimani/whitepointinventory/internal/app"
	"github.com/mike-kimani/whitepointinventory/internal/domain"
	"github.com/mike-kimani/whitepointinventory/pkg/auth"
	"net/http"
)

type UserAuth struct {
	Service app.UserService
}

type authedHandler func(http.ResponseWriter, *http.Request, *domain.User)

func (a *UserAuth) MiddlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			httpresponses.RespondWithError(w, 400, fmt.Sprintf("Auth error: %v", err))
			return
		}
		user, err := a.Service.GetUserByAPIKey(apiKey)

		if err != nil {
			httpresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}
		handler(w, r, user)
	}
}
