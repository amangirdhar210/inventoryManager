package handler

import (
	"encoding/json"
	"net/http"

	"github.com/amangirdhar210/inventory-manager/internal/core/ports"
	"github.com/gorilla/mux"
)

type HTTPHandler struct {
	service ports.InventoryService
}

func NewHTTPHandler(service ports.InventoryService) *HTTPHandler {
	return &HTTPHandler{
		service: service,
	}
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (h *HTTPHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Quantity int     `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	product, err := h.service.AddProduct(req.Name, req.Price, req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	respondWithJSON(w, http.StatusCreated, product)
}

func (h *HTTPHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	product, err := h.service.GetProduct(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	respondWithJSON(w, http.StatusOK, product)
}

func (h *HTTPHandler) SellProductUnits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	_, err := h.service.SellProductUnits(id, req.Quantity)
	if err != nil {
		if err.Error() == "insufficient stock" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		return
	}
	respondWithJSON(w, http.StatusOK, nil)
}

func (h *HTTPHandler) RestockProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		Quantity int `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	_, err := h.service.RestockProduct(id, req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	respondWithJSON(w, http.StatusOK, nil)
}

func (h *HTTPHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.service.DeleteProduct(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	respondWithJSON(w, http.StatusOK, nil)
}

func (h *HTTPHandler) UpdateProductPrice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		NewPrice float64 `json:"price"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateProductPrice(id, req.NewPrice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	respondWithJSON(w, http.StatusOK, nil)
}

func (h *HTTPHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAllProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, products)
}

func (h *HTTPHandler) GetInventoryValue(w http.ResponseWriter, r *http.Request) {
	value, err := h.service.GetInventoryValue()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]float64{"inventory_value": value})
}
