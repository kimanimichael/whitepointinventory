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

type PaymentsHandler struct {
	service     app.PaymentsService
	userService app.UserService
}

func NewPaymentsHandler(service app.PaymentsService) *PaymentsHandler {
	return &PaymentsHandler{
		service: service,
	}
}

func (h *PaymentsHandler) RegisterRoutes(router chi.Router) {
	paymentAuth := middleware.UserAuth{
		Service: h.userService,
	}
	router.Post("/farmers", paymentAuth.MiddlewareAuth(h.CreatePayment))
}

func (h *PaymentsHandler) CreatePayment(w http.ResponseWriter, r *http.Request, user *domain.User) {
	params := CreatePaymentRequest{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		httpresponses.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to decode request body"))
		return
	}
	if params.FarmerName == "" {
		httpresponses.RespondWithError(w, 400, "Farmer name is required")
		return
	}
	if params.CashPaid == 0 || params.ChickenPrice == 0 {
		httpresponses.RespondWithError(w, http.StatusBadRequest, "Both Cash paid and chicken price are required")
		return
	}

	payment, err := h.service.CreatePayment(params.CashPaid, params.ChickenPrice, params.FarmerName, user)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponses.RespondWithJson(w, http.StatusCreated, payment)
}

func (h *PaymentsHandler) GetPaymentByID(w http.ResponseWriter, r *http.Request) {
	params := GetPaymentRequest{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		httpresponses.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to decode request body"))
		return
	}

	payment, err := h.service.GetPaymentByID(params.ID)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponses.RespondWithJson(w, http.StatusOK, payment)
}

func (h *PaymentsHandler) GetPayments(w http.ResponseWriter, r *http.Request) {
	payments, err := h.service.GetPayments()
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponses.RespondWithJson(w, http.StatusOK, payments)
}

func (h *PaymentsHandler) DeletePayment(w http.ResponseWriter, r *http.Request) {
	paymentIDStr := chi.URLParam(r, "payment_id")
	paymentID, err := uuid.Parse(paymentIDStr)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not parse uuid: %s", paymentIDStr))
		return
	}
	err = h.service.DeletePaymentByID(paymentID)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponses.RespondWithJson(w, http.StatusNoContent, nil)
}
