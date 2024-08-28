package purchases

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

// MinChickenNumber prevents zero value chicken number entries
const MinChickenNumber = 1

// MaxChickenNumber prevents entry errors e.g. 1000 instead of 100
const MaxChickenNumber = 999

// MinChickenPrice prevents entry errors e.g. swapping of chicken number and price
const MinChickenPrice = 100

// MaxChickenPrice prevents entry errors e.g. 4500 instead of 450 for price
const MaxChickenPrice = 1000

type ApiConfig struct {
	DB *database.Queries
}

func (apiCfg *ApiConfig) HandlerCreatePurchases(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Chicken    int32  `json:"chicken_no"`
		Price      int32  `json:"chicken_price"`
		FarmerName string `json:"farmer_name"`
	}

	chickenBought := sql.NullFloat64{}
	cashBalance := sql.NullInt32{}

	params := parameters{}

	decode := json.NewDecoder(r.Body)
	err := decode.Decode(&params)
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Error decoding jsonresponses: %v", err))
		return
	}
	if params.FarmerName == "" {
		jsonresponses.RespondWithError(w, 400, "Farmer name is required")
		return
	}
	if params.Chicken < MinChickenNumber || params.Chicken > MaxChickenNumber {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Chicken number must be within %d and %d", MinChickenNumber, MaxChickenNumber))
		return
	}

	if params.Price < MinChickenPrice || params.Price > MaxChickenPrice {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Chicken price must be within %d and %d", MinChickenPrice, MaxChickenPrice))
		return
	}
	farmer, err := apiCfg.DB.GetFarmerByName(r.Context(), params.FarmerName)
	if err != nil {
		jsonresponses.RespondWithError(w, 500, fmt.Sprintf("Error getting farmer: %v", err))
		return
	}

	purchase, err := apiCfg.DB.CreatePurchase(r.Context(), database.CreatePurchaseParams{
		ID:              uuid.New(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Chicken:         params.Chicken,
		PricePerChicken: params.Price,
		UserID:          user.ID,
		FarmerID:        farmer.ID,
	})

	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Error creating purchase: %v", err))
		return
	}

	chickenBought.Float64 = float64(params.Chicken)
	chickenBought.Valid = true

	cashBalance.Int32 = params.Chicken * params.Price
	cashBalance.Valid = true

	err = apiCfg.DB.IncreaseChickenOwed(r.Context(), database.IncreaseChickenOwedParams{
		ChickenBalance: chickenBought,
		ID:             farmer.ID,
	})

	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Error adding chicken owed to farmer database: %v", err))
		return
	}

	err = apiCfg.DB.IncreaseCashOwed(r.Context(), database.IncreaseCashOwedParams{
		CashBalance: cashBalance,
		ID:          farmer.ID,
	})

	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Error adding cash owed to farmer database: %v", err))
		return
	}

	err = apiCfg.DB.MarkFarmerAsUpdated(r.Context(), farmer.ID)
	if err != nil {
		jsonresponses.RespondWithError(w, 500, fmt.Sprintf("Error updating farmer: %v", err))
	}

	fmt.Printf("Farmer %v updated at updated to %v\n", farmer.Name, time.Now())

	jsonresponses.RespondWithJSON(w, 201, purchase)
}

func (apiCfg *ApiConfig) HandlerGetPurchases(w http.ResponseWriter, r *http.Request) {
	purchases, err := apiCfg.DB.GetPurchases(r.Context())
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't fetch purchases: %v", err))
		return
	}
	purchasesWithNames := []models.Purchase{}
	purchasesToPublish := []models.Purchase{}

	for _, purchase := range purchases {
		purchasesWithNames = append(purchasesWithNames, models.DatabasePurchaseToPurchase(purchase))
	}
	for _, purchaseWithName := range purchasesWithNames {
		user, err := apiCfg.DB.GetUserByID(r.Context(), purchaseWithName.UserID)
		if err != nil {
			jsonresponses.RespondWithError(w, 404, fmt.Sprintf("Couldn't fetch user: %v", err))
			return
		}
		purchaseWithName.UserName = user.Name
		farmer, err := apiCfg.DB.GetFarmerByID(r.Context(), purchaseWithName.FarmerID)
		if err != nil {
			jsonresponses.RespondWithError(w, 404, fmt.Sprintf("Couldn't fetch farmer: %v", err))
			return
		}
		purchaseWithName.FarmerName = farmer.Name
		purchaseWithName.FarmerChickenBalance = farmer.ChickenBalance.Float64
		purchaseWithName.FarmerCashBalance = farmer.CashBalance.Int32

		purchasesToPublish = append(purchasesToPublish, purchaseWithName)
	}
	jsonresponses.RespondWithJSON(w, 200, purchasesToPublish)
}

func (apiCfg *ApiConfig) HandlerGetPurchaseByID(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ID uuid.UUID `json:"purchase_id"`
	}
	params := parameters{}

	decode := json.NewDecoder(r.Body)
	err := decode.Decode(&params)
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't decode JSON: %v", err))
		return
	}

	purchase, err := apiCfg.DB.GetPurchaseByID(r.Context(), params.ID)
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't get purchase: %v", err))
	}
	jsonresponses.RespondWithJSON(w, 200, purchase)
}

func (apiCfg *ApiConfig) HandlerDeletePurchase(w http.ResponseWriter, r *http.Request, user database.User) {
	purchaseIDStr := chi.URLParam(r, "purchase_id")
	purchaseID, err := uuid.Parse(purchaseIDStr)
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't parse purchase ID for deletion: %v", err))
		return
	}

	purchase, err := apiCfg.DB.GetPurchaseByID(r.Context(), purchaseID)
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't fetch purchase for deletion: %v", err))
		return
	}

	chickenBalance := sql.NullFloat64{}
	chickenBalance.Float64 = float64(purchase.Chicken)
	chickenBalance.Valid = true

	err = apiCfg.DB.DecreaseChickenOwed(r.Context(), database.DecreaseChickenOwedParams{
		ChickenBalance: chickenBalance,
		ID:             purchase.FarmerID,
	})
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't decrease chicken owed: %v", err))
		return
	}

	cashBalance := sql.NullInt32{}
	cashBalance.Int32 = purchase.Chicken * purchase.PricePerChicken
	cashBalance.Valid = true

	err = apiCfg.DB.DecreaseCashOwed(r.Context(), database.DecreaseCashOwedParams{
		CashBalance: cashBalance,
		ID:          purchase.FarmerID,
	})
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't decrease cash owed: %v", err))
		return
	}

	err = apiCfg.DB.DeletePurchase(r.Context(), purchaseID)
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't parse purchase ID for deletion: %v", err))
		return
	}
	jsonresponses.RespondWithJSON(w, 200, fmt.Sprintf("Purchase successfully deleted by %v", user.Name))

}
