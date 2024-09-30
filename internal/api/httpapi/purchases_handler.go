package httpapi

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
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
