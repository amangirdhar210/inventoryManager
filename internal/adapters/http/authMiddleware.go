package http

import (
	"net/http"
	"strings"

	"github.com/amangirdhar210/inventory-manager/config"
	"github.com/amangirdhar210/inventory-manager/internal/core/domain"
	"github.com/golang-jwt/jwt/v5"
)

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
