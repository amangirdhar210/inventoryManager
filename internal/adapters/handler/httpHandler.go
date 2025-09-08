package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/amangirdhar210/inventory-manager/config"
	"github.com/amangirdhar210/inventory-manager/internal/core/domain"
	"github.com/amangirdhar210/inventory-manager/internal/core/service"
	"github.com/golang-jwt/jwt/v5"
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

func (h *HTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *HTTPHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "logout successful"})
}

func (h *HTTPHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			h.handleError(w, domain.ErrUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			h.handleError(w, domain.ErrUnauthorized)
			return
		}

		claims := &jwt.RegisteredClaims{}
		_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JWTSecretKey), nil
		})

		if err != nil {
			h.handleError(w, domain.ErrUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *HTTPHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Quantity int     `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	product, err := h.inventoryService.AddProduct(req.Name, req.Price, req.Quantity)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusCreated, product)
}

func (h *HTTPHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	product, err := h.inventoryService.GetProduct(id)
	if err != nil {
		h.handleError(w, err)
		return
	}
	h.respondWithJSON(w, http.StatusOK, product)
}

func (h *HTTPHandler) SellProductUnits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	product, err := h.inventoryService.SellProductUnits(id, req.Quantity)
	if err != nil {
		h.handleError(w, err)
		return
	}
	h.respondWithJSON(w, http.StatusOK, product)
}

func (h *HTTPHandler) RestockProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	product, err := h.inventoryService.RestockProduct(id, req.Quantity)
	if err != nil {
		h.handleError(w, err)
		return
	}
	h.respondWithJSON(w, http.StatusOK, product)
}

func (h *HTTPHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.inventoryService.DeleteProduct(id)
	if err != nil {
		h.handleError(w, err)
		return
	}
	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "product deleted successfully"})
}

func (h *HTTPHandler) UpdateProductPrice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var req struct {
		NewPrice float64 `json:"price"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.inventoryService.UpdateProductPrice(id, req.NewPrice)
	if err != nil {
		h.handleError(w, err)
		return
	}
	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "product price updated successfully"})
}

func (h *HTTPHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.inventoryService.GetAllProducts()
	if err != nil {
		h.handleError(w, err)
		return
	}
	h.respondWithJSON(w, http.StatusOK, products)
}

func (h *HTTPHandler) GetInventoryValue(w http.ResponseWriter, r *http.Request) {
	value, err := h.inventoryService.GetInventoryValue()
	if err != nil {
		h.handleError(w, err)
		return
	}
	h.respondWithJSON(w, http.StatusOK, map[string]float64{"inventory_value": value})
}

func (h *HTTPHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (h *HTTPHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

func (h *HTTPHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrProductNotFound):
		h.respondWithError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrInsufficientStock), errors.Is(err, domain.ErrProductInvalid):
		h.respondWithError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrInvalidCredentials), errors.Is(err, domain.ErrUnauthorized):
		h.respondWithError(w, http.StatusUnauthorized, err.Error())
	default:
		h.respondWithError(w, http.StatusInternalServerError, "An internal server error occurred")
	}
}
