package payments

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

// MinCashPaid prevents swapping errors of cash paid and chicken price
const MinCashPaid = 1000

// MaxCashPaid prevents entry errors e.g. 400000 entry instead of 40000
const MaxCashPaid = 100000

// MinChickenPrice prevents entry errors e.g. swapping of chicken number and price
const MinChickenPrice = 100

// MaxChickenPrice prevents entry errors e.g. 4500 instead of 450 for price
const MaxChickenPrice = 1000

type ApiConfig struct {
	DB *database.Queries
}

func (apiCfg *ApiConfig) HandlerCreatePayment(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Cash            int32  `json:"cash_paid"`
		PricePerChicken int32  `json:"price_per_chicken_paid"`
		FarmerName      string `json:"farmer_name"`
	}
	params := parameters{}

	cash_balance := sql.NullInt32{}

	chicken_balance := sql.NullFloat64{}

	decode := json.NewDecoder(r.Body)

	err := decode.Decode(&params)
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't decode jsonresponses data into params: %v", err))
		return
	}
	if params.FarmerName == "" {
		jsonresponses.RespondWithError(w, 400, "Farmer name is required")
		return
	}
	if params.Cash < MinCashPaid || params.Cash > MaxCashPaid {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Cash paid must be between %d and %d", MinCashPaid, MaxCashPaid))
		return
	}
	if params.PricePerChicken < MinChickenPrice || params.PricePerChicken > MaxChickenPrice {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Chicken price must be within %d and %d", MinChickenPrice, MaxChickenPrice))
		return
	}

	farmer, err := apiCfg.DB.GetFarmerByName(r.Context(), params.FarmerName)
	if err != nil {
		jsonresponses.RespondWithError(w, 404, fmt.Sprintf("Couldn't find farmer by name: %v", err))
		return
	}

	payment, err := apiCfg.DB.CreatePayment(r.Context(), database.CreatePaymentParams{
		ID:                  uuid.New(),
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
		CashPaid:            params.Cash,
		PricePerChickenPaid: params.PricePerChicken,
		UserID:              user.ID,
		FarmerID:            farmer.ID,
	})
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't create purchase: %v", err))
		return
	}

	cash_balance.Int32 = params.Cash
	cash_balance.Valid = true
	// TODO: Handle operations that result in floats
	chicken_balance.Float64 = float64(params.Cash) / float64(params.PricePerChicken)
	chicken_balance.Valid = true

	err = apiCfg.DB.DecreaseCashOwed(r.Context(), database.DecreaseCashOwedParams{
		CashBalance: cash_balance,
		ID:          farmer.ID,
	})

	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't decrease cash owed to the farmer: %v", err))
	}

	err = apiCfg.DB.DecreaseChickenOwed(r.Context(), database.DecreaseChickenOwedParams{
		ChickenBalance: chicken_balance,
		ID:             farmer.ID,
	})
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't decrease chicken owed to the farmer: %v", err))
	}
	err = apiCfg.DB.MarkFarmerAsUpdated(r.Context(), farmer.ID)
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't mark farmer as updated: %v", err))
	}

	fmt.Printf("Farmer %v updated at updated to %v\n", farmer.Name, time.Now())

	jsonresponses.RespondWithJSON(w, 201, payment)
}

func (apiCfg *ApiConfig) HandlerGetPayments(w http.ResponseWriter, r *http.Request) {
	customPayments := []models.Payment{}
	paymentResponse := []models.Payment{}

	payments, err := apiCfg.DB.GetPayments(r.Context())
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't fetch payments: %v", err))
	}

	for _, payment := range payments {
		customPayments = append(customPayments, models.DatabasePaymentToPayment(payment))
	}
	for _, customPayment := range customPayments {
		user, err := apiCfg.DB.GetUserByID(r.Context(), customPayment.UserID)
		if err != nil {
			jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't fetch user: %v", err))
			return
		}
		farmer, err := apiCfg.DB.GetFarmerByID(r.Context(), customPayment.FarmerID)
		if err != nil {
			jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't fetch farmer: %v", err))
			return
		}
		customPayment.UserName = user.Name
		customPayment.FarmerName = farmer.Name
		customPayment.FarmerChickenBalance = farmer.ChickenBalance.Float64
		customPayment.FarmerCashBalance = farmer.CashBalance.Int32
		paymentResponse = append(paymentResponse, customPayment)
	}

	jsonresponses.RespondWithJSON(w, 200, paymentResponse)
}

func (apiCfg *ApiConfig) HandlerGetPaymentByID(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ID uuid.UUID `json:"payment_id"`
	}
	params := parameters{}

	decode := json.NewDecoder(r.Body)
	err := decode.Decode(&params)
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't decode jsonresponses: %v", err))
		return
	}

	payment, err := apiCfg.DB.GetPaymentByID(r.Context(), params.ID)
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't get payment: %v", err))
		return
	}

	jsonresponses.RespondWithJSON(w, 200, payment)
}

func (apiCfg *ApiConfig) HandlerDeletePayment(w http.ResponseWriter, r *http.Request, user database.User) {
	paymentIDStr := chi.URLParam(r, "payment_id")
	paymentID, err := uuid.Parse(paymentIDStr)
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't parse payment UUID: %v", err))
	}

	cash_balance := sql.NullInt32{}
	chicken_owed := sql.NullFloat64{}

	payment, err := apiCfg.DB.GetPaymentByID(r.Context(), paymentID)
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't get payment: %v", err))
		return
	}
	cash_balance.Int32 = payment.CashPaid
	cash_balance.Valid = true
	// TODO: Handle operations that result in floats
	chicken_owed.Float64 = float64(payment.CashPaid) / float64(payment.PricePerChickenPaid)
	chicken_owed.Valid = true

	err = apiCfg.DB.IncreaseCashOwed(r.Context(), database.IncreaseCashOwedParams{
		CashBalance: cash_balance,
		ID:          payment.FarmerID,
	})
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't increase cash owned: %v", err))
		return
	}

	err = apiCfg.DB.IncreaseChickenOwed(r.Context(), database.IncreaseChickenOwedParams{
		ChickenBalance: chicken_owed,
		ID:             payment.FarmerID,
	})
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't increase chicken owed: %v", err))
		return
	}

	err = apiCfg.DB.DeletePayments(r.Context(), payment.ID)
	if err != nil {
		jsonresponses.RespondWithError(w, 400, fmt.Sprintf("Couldn't delete payment: %v", err))
		return
	}

	jsonresponses.RespondWithJSON(w, 200, fmt.Sprintf("Deletion successfully done by %v", user.Name))

}
