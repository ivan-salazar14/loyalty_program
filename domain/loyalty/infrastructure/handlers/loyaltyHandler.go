package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/ivan-salazar14/firstGoPackage/domain/loyalty/domain/constants"
	"github.com/ivan-salazar14/firstGoPackage/domain/loyalty/domain/model"
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
		loyaltyService: service.NewLoyaltyService(persistence.NewConnection(db), persistence.NewConnectionQuery(db)),
		router:         router,
	}
}
func (s *LoyaltyHandler) Routes() {
	s.router.Get("/loyalty/{id}", s.GetPointsHandler)
	s.router.Post("/loyalty/redeem", s.RedeemHandler)
	s.router.Post("/loyalty/collect", s.CollectPointsHandler)
	s.router.Get("/loyalty/transactions/{id}", s.GetTransactionsHandler)

}

func (s *LoyaltyHandler) CollectPointsHandler(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	var requestBody struct {
		UserId  string        `json:"userId"`
		Product model.Product `json:"product"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, fmt.Sprintf("failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}
	fmt.Errorf(" requestBody", requestBody)

	err := s.loyaltyService.CollectPoints(ctx, requestBody.UserId, &requestBody.Product)
	if err != nil {
		if errors.Is(err, constants.ErrUserNotFound) {
			http.Error(w, fmt.Sprintf("user with UserID %s not found", requestBody.UserId), http.StatusNotFound)
		} else if errors.Is(err, constants.ErrPointsInsuficient) {
			http.Error(w, fmt.Sprintf("user with UserID %s has insufficient points", requestBody.UserId), http.StatusBadRequest)
		} else {
			http.Error(w, fmt.Sprintf("requestBody: %v", requestBody), http.StatusInternalServerError)
			http.Error(w, fmt.Sprintf("failed to redeem points: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *LoyaltyHandler) RedeemHandler(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	var requestBody struct {
		UserId string `json:"userId"`
		Points int    `json:"points"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, fmt.Sprintf("failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	err := s.loyaltyService.RedeemPoints(ctx, requestBody.UserId, requestBody.Points)
	if err != nil {
		if errors.Is(err, constants.ErrUserNotFound) {
			http.Error(w, fmt.Sprintf("user with UserID %s not found", requestBody.UserId), http.StatusNotFound)
		} else if errors.Is(err, constants.ErrPointsInsuficient) {
			http.Error(w, fmt.Sprintf("user with UserID %s has insufficient points", requestBody.UserId), http.StatusBadRequest)
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
	points, err := s.loyaltyService.LoyaltyQueryRepository.GetPoints(ctx, id)
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

func (s *LoyaltyHandler) GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	var id = chi.URLParam(r, "id")

	transactions, err := s.loyaltyService.GetTransactions(ctx, id)
	if err != nil {
		fmt.Printf("err err err %s\n", err)
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Transactions *[]model.Transaction `json:"transactions"`
	}{
		Transactions: transactions,
	})
}
