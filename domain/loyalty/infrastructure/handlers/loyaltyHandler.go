package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/ivan-salazar14/firstGoPackage/domain/loyalty/domain/constants"
	"github.com/ivan-salazar14/firstGoPackage/domain/loyalty/domain/service"
	"github.com/ivan-salazar14/firstGoPackage/domain/loyalty/infrastructure/persistence"
	"github.com/ivan-salazar14/firstGoPackage/infrastructure/database"
	"net/http"
)

type LoyaltyHandler struct {
	loyaltyService *service.LoyaltyService
	router         *chi.Mux
}

func NewLoyaltyHandler(db *database.Service, router *chi.Mux) *LoyaltyHandler {
	return &LoyaltyHandler{
		loyaltyService: service.NewLoyaltyService(persistence.NewConnection(db)),
		router:         router,
	}
}
func (s *LoyaltyHandler) Routes() {
	s.router.Get("/loyalty/{id}", s.GetPointsHandler)
	s.router.Post("/loyalty/redeem/{id}", s.RedeemHandler)
	s.router.Post("/loyalty/collect/{id}", s.RedeemHandler)

}

func (s *LoyaltyHandler) CollectPointsHandler(w http.ResponseWriter, r *http.Request) {
	var id = chi.URLParam(r, "id")
	var ctx = r.Context()

	var p string
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, fmt.Sprintf("failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	err := s.loyaltyService.CollectPoints(ctx, id, p)
	if err != nil {
		if errors.Is(err, constants.ErrUserNotFound) {
			http.Error(w, fmt.Sprintf("user with ID %s not found", id), http.StatusNotFound)
		} else if errors.Is(err, constants.ErrPointsInsuficient) {
			http.Error(w, fmt.Sprintf("user with ID %s has insufficient points", id), http.StatusBadRequest)
		} else {
			http.Error(w, fmt.Sprintf("failed to redeem points: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *LoyaltyHandler) RedeemHandler(w http.ResponseWriter, r *http.Request) {
	var id = chi.URLParam(r, "id")
	var ctx = r.Context()

	var p int
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, fmt.Sprintf("failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	err := s.loyaltyService.RedeemPoints(ctx, id, p)
	if err != nil {
		if errors.Is(err, constants.ErrUserNotFound) {
			http.Error(w, fmt.Sprintf("user with ID %s not found", id), http.StatusNotFound)
		} else if errors.Is(err, constants.ErrPointsInsuficient) {
			http.Error(w, fmt.Sprintf("user with ID %s has insufficient points", id), http.StatusBadRequest)
		} else {
			http.Error(w, fmt.Sprintf("failed to redeem points: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *LoyaltyHandler) GetPointsHandler(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	var id = chi.URLParam(r, "id")
	fmt.Printf("id %s\n", id)
	points, err := s.loyaltyService.GetPoints(ctx, id)
	if err != nil {
		fmt.Printf("err err err %s\n", err)
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Points int `json:"points"`
	}{
		Points: points,
	})
}
