package ports

import "github.com/amangirdhar210/inventory-manager/internal/core/domain"

type ProductRepository interface {
	FindById(id string) (*domain.Product, error)
	ListAll() ([]domain.Product, error)
	Save(product *domain.Product) error
	Update(product *domain.Product) error
	DeleteById(id string) error
}

type Notifier interface {
	NotifyLowStock(product *domain.Product)
}
