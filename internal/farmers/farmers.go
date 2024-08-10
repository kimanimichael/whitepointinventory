package farmers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/mike-kimani/whitepointinventory/internal/models"
	"github.com/mike-kimani/whitepointinventory/pkg/jsonresponses"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/mike-kimani/whitepointinventory/internal/database"
)

type ApiConfig struct {
	DB *database.Queries
}

func (apiCfg *ApiConfig) HandlerCreateFarmer(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name           string  `json:"name"`
		ChickenBalance float64 `json:"chicken_balance"`
		CashBalance    int32   `json:"cash_balance"`
	}

	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't parse jsonresponses %v", err))
		return
	}

	chickenBalance := sql.NullFloat64{}

	if params.ChickenBalance != 0 {
		chickenBalance.Float64 = params.ChickenBalance
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
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't create farmer: %v", err))
		return
	}

	jsonresponses.RespondWithJSON(w, 200, farmer)
}

func (apiCfg *ApiConfig) HandlerGetFarmerByName(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	params := parameters{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't parse jsonresponses: %v", err))
		return
	}

	farmer, err := apiCfg.DB.GetFarmerByName(r.Context(), params.Name)
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't get farmer: %v", err))
		return
	}

	jsonresponses.RespondWithJSON(w, 200, farmer)
}

func (apiCfg *ApiConfig) HandlerDeleteFarmer(w http.ResponseWriter, r *http.Request) {

	farmerIDStr := chi.URLParam(r, "farmer_id")
	farmerID, err := uuid.Parse(farmerIDStr)
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't parse UUID for deletion: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFarmers(r.Context(), farmerID)
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't delete farmer: %v", err))
		return
	}
	jsonresponses.RespondWithJSON(w, 200, struct{}{})
}

func (apiCfg *ApiConfig) HandlerGetFarmers(w http.ResponseWriter, r *http.Request) {
	farmers, err := apiCfg.DB.GetFarmers(r.Context())
	if err != nil {
		jsonresponses.RespondWithError(w, 404, fmt.Sprintf("Couldn't get farmers: %v", err))
	}
	jsonresponses.RespondWithJSON(w, 200, models.DatabaseFarmersToFarmers(farmers))
}
