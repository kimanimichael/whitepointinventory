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
		password, name, err := auth.GetPasswordAndName(r.Header)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByPasswordAndUserName(r.Context(), database.GetUserByPasswordAndUserNameParams{
			Password: password,
			Name: name,
		})
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}
		handler(w, r, user)
	}
}