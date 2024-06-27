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

func (apiCfg *apiConfig) handerCreatePurchases(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Chicken    int32  `json:"chicken_no"`
		Price      int32  `json:"chicken_price"`
		FarmerName string `json:"farmer_name"`
	}

	chicken_bought := sql.NullInt32{}
	cash_balance := sql.NullInt32{}

	params := parameters{}

	decode := json.NewDecoder(r.Body)
	err := decode.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error decoding json: %v", err))
		return
	}
	farmer, err := apiCfg.DB.GetFarmerByName(r.Context(), params.FarmerName)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error getting farmer: %v", err))
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
		respondWithError(w, 400, fmt.Sprintf("Error creating purchase: %v", err))
		return
	}

	chicken_bought.Int32 = params.Chicken
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

	chicken_balance := sql.NullInt32{}
	chicken_balance.Int32 = purchase.Chicken
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
