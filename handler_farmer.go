package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/mike-kimani/whitepointinventory/internal/database"
)

func (apiCfg *apiConfig) handlerCreateFarmer(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name           string `json:"name"`
		ChickenBalance int32  `json:"chicken_balance"`
		CashBalance    int32  `json:"cash_balance"`
	}

	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse json %v", err))
		return
	}

	chickenBalance := sql.NullInt32{}

	if params.ChickenBalance != 0 {
		chickenBalance.Int32 = params.ChickenBalance
		chickenBalance.Valid = true
	}

	cashBalance := sql.NullInt32{}

	if params.CashBalance != 0 {
		cashBalance.Int32 = params.CashBalance
		cashBalance.Valid = true
	}

	farmer, err := apiCfg.DB.CreateFarmer(r.Context(), database.CreateFarmerParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Name:           params.Name,
		ChickenBalance: chickenBalance,
		CashBalance:    cashBalance,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create farmer: %v", err))
		return
	}

	respondWithJSON(w, 200, farmer)
}

func (apiCfg *apiConfig) handlerGetFarmerByName(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	params := parameters{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse json: %v", err))
		return
	}

	farmer, err := apiCfg.DB.GetFarmerByName(r.Context(), params.Name)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get farmer: %v", err))
		return
	}

	respondWithJSON(w, 200, farmer)
}

func (apiCfg *apiConfig) handlerDeleteFarmer(w http.ResponseWriter, r *http.Request) {

	farmerIDstr := chi.URLParam(r, "farmer_id")
	farmerID, err := uuid.Parse(farmerIDstr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse UUID for deletion: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFarmers(r.Context(), farmerID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't delete farmer: %v", err))
		return
	}
	respondWithJSON(w, 200, struct{}{})
}

func (apiCfg *apiConfig) handlerGetFarmers(w http.ResponseWriter, r *http.Request) {
	farmers, err := apiCfg.DB.GetFarmers(r.Context())
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("Couldn't get farmers: %v", err))
	}
	respondWithJSON(w, 200, databaseFarmersToFarmers(farmers))
}
