package farmersapi

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/mike-kimani/fechronizo/v2/pkg/httpresponses"
	"github.com/mike-kimani/whitepointinventory/internal/farmers"
	httpapi "github.com/mike-kimani/whitepointinventory/internal/http"
	"github.com/mike-kimani/whitepointinventory/internal/users"
	"log"
	"net/http"
)

type FarmerHandler struct {
	service     farmers.FarmerService
	userService users.UserService
}

func NewFarmerHandler(service farmers.FarmerService, userService users.UserService) *FarmerHandler {
	return &FarmerHandler{
		service:     service,
		userService: userService,
	}
}

func (h *FarmerHandler) RegisterRoutes(router chi.Router) {
	farmerAuth := httpapi.UserAuth{
		Service: h.userService,
	}
	router.Post("/farmers", h.CreateFarmer)
	router.Get("/farmers", h.GetFarmerByName)
	router.Get("/farmer", h.GetFarmers)
	router.Get("/paged_farmers", h.GetPagedFarmers)
	router.Post("/set_farmer_balance", farmerAuth.MiddlewareAuth(h.SetFarmerBalances))
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

func (h *FarmerHandler) SetFarmerBalances(w http.ResponseWriter, r *http.Request, user *users.User) {
	params := SetFarmerBalancesRequest{}

	decode := json.NewDecoder(r.Body)
	if err := decode.Decode(&params); err != nil {
		httpresponses.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not decode parameters :%v", err))
	}
	ctx := r.Context()
	updatedFarmer, err := h.service.SetFarmerBalances(ctx, params.Name, params.ChickenBalance, params.CashBalance)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not update farmer :%v", err))
		return
	}
	log.Printf("Farmer %s updated by %s\n", updatedFarmer.Name, user.Name)
	httpresponses.RespondWithJson(w, http.StatusOK, farmerToResponseFarmer(*updatedFarmer))
}

func (h *FarmerHandler) DeleteFarmerByID(w http.ResponseWriter, r *http.Request) {
	farmerID := chi.URLParam(r, "farmerID")

	ctx := r.Context()

	err := h.service.DeleteFarmerByID(ctx, farmerID)
	if err != nil {
		httpresponses.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not delete farmer :%v", err))
		return
	}
	httpresponses.RespondWithJson(w, http.StatusNoContent, nil)
}
