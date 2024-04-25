package main

import (
	"fmt"
	"net/http"

	"github.com/mike-kimani/whitepointinventory/auth"
	"github.com/mike-kimani/whitepointinventory/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth (handler authedHandler) http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request)  {
		password, email, err := auth.GetPasswordAndEmail(r.Header)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByPasswordAndEmail(r.Context(), database.GetUserByPasswordAndEmailParams{
			Password: password,
			Email: email,
		})
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}
		handler(w, r, user)
	}
}