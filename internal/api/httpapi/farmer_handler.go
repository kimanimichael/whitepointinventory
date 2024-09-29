package httpapi

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/mike-kimani/fechronizo/v2/pkg/httpresponses"
	"github.com/mike-kimani/whitepointinventory/internal/app"
	"net/http"
)

type FarmerHandler struct {
	service app.FarmerService
}

func NewFarmerHandler(service app.FarmerService) *FarmerHandler {
	return &FarmerHandler{
		service: service,
	}
}

func (h *FarmerHandler) RegisterRoutes(router chi.Router) {
	router.Post("/farmer", h.CreateFarmer)
}

func (h *FarmerHandler) CreateFarmer(w http.ResponseWriter, r *http.Request) {
	params := CreateFarmerRequest{}

	decode := json.NewDecoder(r.Body)
	if err := decode.Decode(&params); err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not decode parameters :%v", err))
		return
	}

	farmer, err := h.service.CreateFarmer(params.Name, params.ChickenBalance, params.CashBalance)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not create farmer :%v", err))
		return
	}
	httpresponses.RespondWithJson(w, http.StatusCreated, farmer)
}
