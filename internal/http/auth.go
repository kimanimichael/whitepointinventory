package httpapi

import (
	"fmt"
	"github.com/mike-kimani/fechronizo/v2/pkg/httpresponses"
	"github.com/mike-kimani/whitepointinventory/internal/users"
	"github.com/mike-kimani/whitepointinventory/pkg/auth"
	"net/http"
)

type UserAuth struct {
	Service users.UserService
}

type authedHandler func(http.ResponseWriter, *http.Request, *users.User)

func (a *UserAuth) MiddlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			httpresponses.RespondWithError(w, 400, fmt.Sprintf("Auth error: %v", err))
			return
		}
		ctx := r.Context()
		user, err := a.Service.GetUserByAPIKey(ctx, apiKey)

		if err != nil {
			httpresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}
		handler(w, r, user)
	}
}
