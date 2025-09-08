package ports

import "github.com/amangirdhar210/inventory-manager/internal/core/domain"

type ManagerRepository interface {
	FindByEmail(email string) (*domain.Manager, error)
}
