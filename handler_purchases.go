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

// MinChickenNumber prevents zero value chicken number entries
const MinChickenNumber = 1

// MaxChickenNumber prevents entry errors e.g. 1000 instead of 100
const MaxChickenNumber = 999

// MinChickenPrice prevents entry errors e.g. swapping of chicken number and price
const MinChickenPrice = 100

// MaxChickenPrice prevents entry errors e.g. 4500 instead of 450 for price
const MaxChickenPrice = 1000

func (apiCfg *apiConfig) handerCreatePurchases(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Chicken    int32  `json:"chicken_no"`
		Price      int32  `json:"chicken_price"`
		FarmerName string `json:"farmer_name"`
	}

	chicken_bought := sql.NullFloat64{}
	cash_balance := sql.NullInt32{}

	params := parameters{}

	decode := json.NewDecoder(r.Body)
	err := decode.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error decoding json: %v", err))
		return
	}
	if params.FarmerName == "" {
		respondWithError(w, 400, "Farmer name is required")
		return
	}
	if params.Chicken < MinChickenNumber || params.Chicken > MaxChickenNumber {
		respondWithError(w, 400, fmt.Sprintf("Chicken number must be within %d and %d", MinChickenNumber, MaxChickenNumber))
		return
	}

	if params.Price < MinChickenPrice || params.Price > MaxChickenPrice {
		respondWithError(w, 400, fmt.Sprintf("Chicken price must be within %d and %d", MinChickenPrice, MaxChickenPrice))
		return
	}
	farmer, err := apiCfg.DB.GetFarmerByName(r.Context(), params.FarmerName)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error getting farmer: %v", err))
		return
	}
	mostRecentPurchase, err := apiCfg.DB.GetMostRecentPurchase(r.Context())
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error getting most recent purchase: %v", err))
	}
	currentTime := time.Now()
	correctedRecentPurchaseTime := time.Date(
		mostRecentPurchase.CreatedAt.Year(),
		mostRecentPurchase.CreatedAt.Month(),
		mostRecentPurchase.CreatedAt.Day(),
		mostRecentPurchase.CreatedAt.Hour(),
		mostRecentPurchase.CreatedAt.Minute(),
		mostRecentPurchase.CreatedAt.Second(),
		mostRecentPurchase.CreatedAt.Nanosecond(),
		time.FixedZone("EAT", 3*60*60),
	)
	durationSinceLastPayment := currentTime.Sub(correctedRecentPurchaseTime)

	if mostRecentPurchase.FarmerID == farmer.ID {
		if mostRecentPurchase.Chicken == params.Chicken {
			if mostRecentPurchase.PricePerChicken == params.Price {
				if durationSinceLastPayment < IdenticalTransactionInterval {
					fmt.Printf("Identical transactions in less than %ds attempted", int(IdenticalTransactionInterval.Seconds()))
					respondWithError(w, 404, fmt.Sprintf("Similar transaction for %s. Wait for %ds", farmer.Name, int(IdenticalTransactionInterval.Seconds()-durationSinceLastPayment.Seconds())))
					return
				}
			}
		}
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
		respondWithError(w, 400, fmt.Sprintf("Error creating purchase: %v", err))
		return
	}

	chicken_bought.Float64 = float64(params.Chicken)
	chicken_bought.Valid = true

	cash_balance.Int32 = params.Chicken * params.Price
	cash_balance.Valid = true

	err = apiCfg.DB.IncreaseChickenOwed(r.Context(), database.IncreaseChickenOwedParams{
		ChickenBalance: chicken_bought,
		ID:             farmer.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error adding chicken owed to farmer database: %v", err))
		return
	}

	err = apiCfg.DB.IncreaseCashOwed(r.Context(), database.IncreaseCashOwedParams{
		CashBalance: cash_balance,
		ID:          farmer.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error adding cash owed to farmer database: %v", err))
		return
	}

	err = apiCfg.DB.MarkFarmerAsUpdated(r.Context(), farmer.ID)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error updating farmer: %v", err))
	}

	fmt.Printf("Farmer %v updated at updated to %v\n", farmer.Name, time.Now())

	respondWithJSON(w, 201, purchase)
}

func (apiCfg *apiConfig) handlerGetPurchases(w http.ResponseWriter, r *http.Request) {
	purchases, err := apiCfg.DB.GetPurchases(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't fetch purchases: %v", err))
		return
	}
	purchasesWithNames := []Purchase{}
	purchasesToPublish := []Purchase{}

	for _, purchase := range purchases {
		purchasesWithNames = append(purchasesWithNames, databasePurchaseToPurchase(purchase))
	}
	for _, purchaseWithName := range purchasesWithNames {
		user, err := apiCfg.DB.GetUserByID(r.Context(), purchaseWithName.UserID)
		if err != nil {
			respondWithError(w, 404, fmt.Sprintf("Couldn't fetch user: %v", err))
			return
		}
		purchaseWithName.UserName = user.Name
		farmer, err := apiCfg.DB.GetFarmerByID(r.Context(), purchaseWithName.FarmerID)
		if err != nil {
			respondWithError(w, 404, fmt.Sprintf("Couldn't fetch farmer: %v", err))
			return
		}
		purchaseWithName.FarmerName = farmer.Name
		purchaseWithName.FarmerChickenBalance = farmer.ChickenBalance.Float64
		purchaseWithName.FarmerCashBalance = farmer.CashBalance.Int32

		purchasesToPublish = append(purchasesToPublish, purchaseWithName)
	}
	respondWithJSON(w, 200, purchasesToPublish)
}

func (apiCfg *apiConfig) handlerGetPurchaseByID(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ID uuid.UUID `json:"purchase_id"`
	}
	params := parameters{}

	decode := json.NewDecoder(r.Body)
	err := decode.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't decode JSON: %v", err))
		return
	}

	purchase, err := apiCfg.DB.GetPurchaseByID(r.Context(), params.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get purchase: %v", err))
	}
	respondWithJSON(w, 200, purchase)
}

func (apiCfg *apiConfig) handlerDeletePurchase(w http.ResponseWriter, r *http.Request, user database.User) {
	purchaseIDStr := chi.URLParam(r, "purchase_id")
	purchaseID, err := uuid.Parse(purchaseIDStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse purchase ID for deletion: %v", err))
		return
	}

	purchase, err := apiCfg.DB.GetPurchaseByID(r.Context(), purchaseID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't fetch purchase for deletion: %v", err))
		return
	}

	chicken_balance := sql.NullFloat64{}
	chicken_balance.Float64 = float64(purchase.Chicken)
	chicken_balance.Valid = true

	err = apiCfg.DB.DecreaseChickenOwed(r.Context(), database.DecreaseChickenOwedParams{
		ChickenBalance: chicken_balance,
		ID:             purchase.FarmerID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't decrease chicken owed: %v", err))
		return
	}

	cash_balance := sql.NullInt32{}
	cash_balance.Int32 = purchase.Chicken * purchase.PricePerChicken
	cash_balance.Valid = true

	err = apiCfg.DB.DecreaseCashOwed(r.Context(), database.DecreaseCashOwedParams{
		CashBalance: cash_balance,
		ID:          purchase.FarmerID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't decrease cash owed: %v", err))
		return
	}

	err = apiCfg.DB.DeletePurchase(r.Context(), purchaseID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse purchase ID for deletion: %v", err))
		return
	}
	respondWithJSON(w, 200, fmt.Sprintf("Purchase successfully deleted by %v", user.Name))

}
