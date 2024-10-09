package httpapi

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
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
	router.Post("/farmers", h.CreateFarmer)
	router.Get("/farmers", h.GetFarmerByName)
	router.Get("/farmer", h.GetFarmers)
	router.Delete("/farmer", h.DeleteFarmerByID)
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

func (h *FarmerHandler) GetFarmerByName(w http.ResponseWriter, r *http.Request) {
	params := GetFarmerRequest{}

	decode := json.NewDecoder(r.Body)
	if err := decode.Decode(&params); err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not decode parameters :%v", err))
		return
	}
	farmer, err := h.service.GetFarmerByName(params.Name)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not get farmer :%v", err))
		return
	}
	httpresponses.RespondWithJson(w, http.StatusOK, farmer)
}

func (h *FarmerHandler) GetFarmers(w http.ResponseWriter, r *http.Request) {
	farmers, err := h.service.GetFarmers()
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not get farmers :%v", err))
		return
	}
	farmersResponse := DomainFarmersToFarmers(farmers)
	httpresponses.RespondWithJson(w, http.StatusOK, farmersResponse)
}

func (h *FarmerHandler) DeleteFarmerByID(w http.ResponseWriter, r *http.Request) {
	farmerIDStr := chi.URLParam(r, "farmerID")
	farmerID, err := uuid.Parse(farmerIDStr)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not parse farmerID :%v", err))
		return
	}
	err = h.service.DeleteFarmerByID(farmerID)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not delete farmer :%v", err))
		return
	}
	httpresponses.RespondWithJson(w, http.StatusNoContent, nil)
}
