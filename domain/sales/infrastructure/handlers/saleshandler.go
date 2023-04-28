package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/ivan-salazar14/firstGoPackage/domain/sales/domain/model"
	"github.com/ivan-salazar14/firstGoPackage/domain/sales/domain/service"
	"github.com/ivan-salazar14/firstGoPackage/domain/sales/infrastructure/persistence"
	"github.com/ivan-salazar14/firstGoPackage/infrastructure/database"
	"net/http"
)

type SaleHandler struct {
	saleService *service.SaleService
	router      *chi.Mux
}

func NewSaleHandler(db *database.DataDB, router *chi.Mux) *SaleHandler {
	return &SaleHandler{
		saleService: service.NewSaleService(persistence.NewConnection(db)),
		router:      router,
	}
}
func (s *SaleHandler) Routes() {
	s.router.Get("/sales/{id}", s.GetSaleHandler)
	s.router.Post("/sales", s.CreateSaleHandler)
}
func (s *SaleHandler) CreateSaleHandler(w http.ResponseWriter, r *http.Request) {
	var sale model.Sale
	var ctx = r.Context()
	err := json.NewDecoder(r.Body).Decode(&sale)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	resp, err := s.saleService.CreateSale(ctx, &sale)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": sale.Id, "message": resp})
}

func (s *SaleHandler) GetSaleHandler(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	var id = chi.URLParam(r, "id")
	fmt.Printf("saleid %s\n", id)
	sale, err := s.saleService.GetSale(ctx, id)
	if err != nil {
		fmt.Printf("err err err %s\n", err)
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sale)
}
