package purchasesapi

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/mike-kimani/fechronizo/v2/pkg/httpresponses"
	"github.com/mike-kimani/whitepointinventory/internal/http"
	"github.com/mike-kimani/whitepointinventory/internal/purchases"
	"github.com/mike-kimani/whitepointinventory/internal/users"
	"net/http"
)

type PurchasesHandler struct {
	service     purchases.PurchaseService
	userService users.UserService
}

func NewPurchasesHandler(service purchases.PurchaseService, userService users.UserService) *PurchasesHandler {
	return &PurchasesHandler{
		service:     service,
		userService: userService,
	}
}

func (h *PurchasesHandler) RegisterRoutes(router chi.Router) {
	purchasesAuth := httpapi.UserAuth{
		Service: h.userService,
	}
	router.Post("/purchases", purchasesAuth.MiddlewareAuth(h.CreatePurchase))
	router.Get("/purchase", h.GetPurchaseByID)
	router.Get("/purchases", h.GetPurchases)
	router.Get("/paged_purchases", h.GetPagedPurchases)
	router.Delete("/purchases/{purchase_id}", purchasesAuth.MiddlewareAuth(h.DeletePurchase))
}

func (h *PurchasesHandler) CreatePurchase(w http.ResponseWriter, r *http.Request, user *users.User) {
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

	ctx := r.Context()

	purchase, err := h.service.CreatePurchase(ctx, params.ChickenNo, params.ChickenPrice, params.FarmerName, user)
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

	ctx := r.Context()

	purchase, err := h.service.GetPurchaseByID(ctx, params.ID)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponses.RespondWithJson(w, http.StatusOK, purchase)
}

func (h *PurchasesHandler) GetPurchases(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fetchedPurchases, err := h.service.GetPurchases(ctx)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	purchasesResponse := purchaseToPurchaseResponses(fetchedPurchases)
	httpresponses.RespondWithJson(w, http.StatusOK, purchasesResponse)
}

func (h *PurchasesHandler) GetPagedPurchases(w http.ResponseWriter, r *http.Request) {
	params := GetPagedPurchasesRequest{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		httpresponses.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to decode request body"))
		return
	}
	ctx := r.Context()
	pagedPurchases, err := h.service.GetPagedPurchases(ctx, params.Offset, params.Limit)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	httpresponses.RespondWithJson(w, http.StatusOK, pagedPurchases)
}

func (h *PurchasesHandler) DeletePurchase(w http.ResponseWriter, r *http.Request, user *users.User) {
	purchaseID := chi.URLParam(r, "purchase_id")

	ctx := r.Context()

	err := h.service.DeletePurchaseByID(ctx, purchaseID)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponses.RespondWithJson(w, http.StatusOK, fmt.Sprintf("Purchase successfully deleted by %v", user.Name))
}
