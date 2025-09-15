package ports

import "github.com/amangirdhar210/inventory-manager/internal/core/domain"

type Notifier interface {
	NotifyLowStock(product *domain.Product)
}
