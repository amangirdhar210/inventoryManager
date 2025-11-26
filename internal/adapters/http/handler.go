package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/amangirdhar210/inventory-manager/internal/core/service"
	"github.com/gorilla/mux"
)

type HTTPHandler struct {
	inventoryService service.InventoryService
	authService      service.AuthService
}

func NewHTTPHandler(invService service.InventoryService, authService service.AuthService) *HTTPHandler {
	return &HTTPHandler{
		inventoryService: invService,
		authService:      authService,
	}
}

func (handler *HTTPHandler) Logout(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "logout successful"})
}

func (handler *HTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	log.Printf("%+v", req)

	token, err := handler.authService.Login(req.Email, req.Password)
	if err != nil {
		handleError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (handler *HTTPHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Quantity int     `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	product, err := handler.inventoryService.AddProduct(req.Name, req.Price, req.Quantity)
	if err != nil {
		handleError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, product)
}

func (handler *HTTPHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	product, err := handler.inventoryService.GetProduct(id)
	if err != nil {
		handleError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, product)
}

func (handler *HTTPHandler) SellProductUnits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	err := handler.inventoryService.SellProductUnits(id, req.Quantity)
	if err != nil {
		handleError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "sell request processed successfully"})
}

func (handler *HTTPHandler) RestockProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	err := handler.inventoryService.RestockProduct(id, req.Quantity)
	if err != nil {
		handleError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "restock request processed successfully."})
}

func (handler *HTTPHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := handler.inventoryService.DeleteProduct(id)
	if err != nil {
		handleError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "product deleted successfully"})
}

func (handler *HTTPHandler) UpdateProductPrice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var req struct {
		NewPrice float64 `json:"price"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := handler.inventoryService.UpdateProductPrice(id, req.NewPrice)
	if err != nil {
		handleError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "product price updated successfully"})
}

func (handler *HTTPHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := handler.inventoryService.GetAllProducts()
	if err != nil {
		handleError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, products)
}

func (handler *HTTPHandler) GetInventoryValue(w http.ResponseWriter, r *http.Request) {
	value, err := handler.inventoryService.GetInventoryValue()
	if err != nil {
		handleError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]float64{"inventory_value": value})
}
