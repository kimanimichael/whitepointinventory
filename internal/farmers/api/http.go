package farmersapi

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/mike-kimani/fechronizo/v2/pkg/httpresponses"
	"github.com/mike-kimani/whitepointinventory/internal/farmers"
	"net/http"
)

type FarmerHandler struct {
	service farmers.FarmerService
}

func NewFarmerHandler(service farmers.FarmerService) *FarmerHandler {
	return &FarmerHandler{
		service: service,
	}
}

func (h *FarmerHandler) RegisterRoutes(router chi.Router) {
	router.Post("/farmers", h.CreateFarmer)
	router.Get("/farmers", h.GetFarmerByName)
	router.Get("/farmer", h.GetFarmers)
	router.Get("/paged_farmers", h.GetPagedFarmers)
	router.Delete("/farmer", h.DeleteFarmerByID)
}

func (h *FarmerHandler) CreateFarmer(w http.ResponseWriter, r *http.Request) {
	params := CreateFarmerRequest{}

	decode := json.NewDecoder(r.Body)
	if err := decode.Decode(&params); err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not decode parameters :%v", err))
		return
	}

	if params.Name == "" {
		httpresponses.RespondWithError(w, http.StatusBadRequest, "Farmer name is required")
		return
	}

	ctx := r.Context()

	farmer, err := h.service.CreateFarmer(ctx, params.Name, params.ChickenBalance, params.CashBalance)
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

	ctx := r.Context()

	farmer, err := h.service.GetFarmerByName(ctx, params.Name)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not get farmer :%v", err))
		return
	}
	httpresponses.RespondWithJson(w, http.StatusOK, farmer)
}

func (h *FarmerHandler) GetFarmers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fetchedFarmers, err := h.service.GetFarmers(ctx)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not get fetchedFarmers :%v", err))
		return
	}
	farmersResponse := farmersToResponseFarmers(fetchedFarmers)
	httpresponses.RespondWithJson(w, http.StatusOK, farmersResponse)
}

func (h *FarmerHandler) GetPagedFarmers(w http.ResponseWriter, r *http.Request) {
	params := GetPagedFarmersRequest{}

	decode := json.NewDecoder(r.Body)
	if err := decode.Decode(&params); err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not decode parameters :%v", err))
		return
	}
	ctx := r.Context()

	fetchedFarmers, err := h.service.GetPagedFarmers(ctx, params.Offset, params.Limit)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not get fetchedFarmers :%v", err))
		return
	}
	httpresponses.RespondWithJson(w, http.StatusOK, fetchedFarmers)
}

func (h *FarmerHandler) DeleteFarmerByID(w http.ResponseWriter, r *http.Request) {
	farmerIDStr := chi.URLParam(r, "farmerID")
	farmerID, err := uuid.Parse(farmerIDStr)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not parse farmerID :%v", err))
		return
	}

	ctx := r.Context()

	err = h.service.DeleteFarmerByID(ctx, farmerID)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not delete farmer :%v", err))
		return
	}
	httpresponses.RespondWithJson(w, http.StatusNoContent, nil)
}
