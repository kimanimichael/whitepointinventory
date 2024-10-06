package httpapi

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/mike-kimani/fechronizo/v2/pkg/httpresponses"
	"github.com/mike-kimani/whitepointinventory/internal/app"
	"github.com/mike-kimani/whitepointinventory/internal/domain"
	"github.com/mike-kimani/whitepointinventory/internal/middleware"
	"net/http"
)

type PurchasesHandler struct {
	service     app.PurchaseService
	userService app.UserService
}

func NewPurchasesHandler(service app.PurchaseService) *PurchasesHandler {
	return &PurchasesHandler{
		service: service,
	}
}

func (h *PurchasesHandler) RegisterRoutes(router chi.Router) {
	purchasesAuth := middleware.UserAuth{
		Service: h.userService,
	}
	router.Post("/purchase", purchasesAuth.MiddlewareAuth(h.CreatePurchase))
	router.Get("/purchase", h.GetPurchaseByID)
	router.Get("/purchases", h.GetPurchases)
	router.Delete("/purchases/{purchase_id}", h.DeletePurchase)
}

func (h *PurchasesHandler) CreatePurchase(w http.ResponseWriter, r *http.Request, user *domain.User) {
	params := CreatePurchaseRequest{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		httpresponses.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Request body could not be decoded as JSON: %v", err))
		return
	}

	if params.FarmerName == "" {
		httpresponses.RespondWithError(w, http.StatusBadRequest, "Farmer name is required")
		return
	}
	if params.ChickenNo == 0 || params.ChickenPrice == 0 {
		httpresponses.RespondWithError(w, http.StatusBadRequest, "Both Chicken number and chicken price are required")
		return
	}

	purchase, err := h.service.CreatePurchase(params.ChickenNo, params.ChickenPrice, params.FarmerName, user)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponses.RespondWithJson(w, http.StatusOK, purchase)
}

func (h *PurchasesHandler) GetPurchaseByID(w http.ResponseWriter, r *http.Request) {
	params := GetTransactionRequest{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		httpresponses.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to decode request body"))
		return
	}
	purchase, err := h.service.GetPurchaseByID(params.ID)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponses.RespondWithJson(w, http.StatusOK, purchase)
}

func (h *PurchasesHandler) GetPurchases(w http.ResponseWriter, r *http.Request) {
	purchases, err := h.service.GetPurchases()
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponses.RespondWithJson(w, http.StatusOK, purchases)
}

func (h *PurchasesHandler) DeletePurchase(w http.ResponseWriter, r *http.Request) {
	purchaseIDStr := chi.URLParam(r, "purchase_id")
	purchaseID, err := uuid.Parse(purchaseIDStr)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not parse uuid: %s", purchaseIDStr))
		return
	}
	err = h.service.DeletePurchaseByID(purchaseID)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponses.RespondWithJson(w, http.StatusNoContent, nil)
}
