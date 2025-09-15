package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/amangirdhar210/inventory-manager/internal/core/domain"
)

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrProductNotFound):
		respondWithError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrInsufficientStock), errors.Is(err, domain.ErrProductInvalid):
		respondWithError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrInvalidCredentials), errors.Is(err, domain.ErrUnauthorized):
		respondWithError(w, http.StatusUnauthorized, err.Error())
	default:
		respondWithError(w, http.StatusInternalServerError, "An internal server error occurred")
	}
}
