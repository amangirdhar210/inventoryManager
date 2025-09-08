package ports

import "github.com/amangirdhar210/inventory-manager/internal/core/domain"

type TokenGenerator interface {
	GenerateToken(manager *domain.Manager) (string, error)
}
