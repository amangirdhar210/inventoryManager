package auth

import (
	"time"

	"github.com/amangirdhar210/inventory-manager/internal/core/domain"
	"github.com/amangirdhar210/inventory-manager/internal/core/ports"
	"github.com/golang-jwt/jwt/v5"
)

var _ ports.TokenGenerator = (*JWTGenerator)(nil)

type JWTGenerator struct {
	secretKey string
}

func NewJWTGenerator(secretKey string) *JWTGenerator {
	return &JWTGenerator{secretKey: secretKey}
}

func (g *JWTGenerator) GenerateToken(manager *domain.Manager) (string, error) {
	claims := &jwt.RegisteredClaims{
		Issuer:    "inventory-manager",
		Subject:   manager.Id,
		Audience:  jwt.ClaimStrings{"managers"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(g.secretKey))
}
