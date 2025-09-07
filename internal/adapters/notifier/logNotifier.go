package notifier

import (
	"log"

	"github.com/amangirdhar210/inventory-manager/internal/core/domain"
	"github.com/amangirdhar210/inventory-manager/internal/core/ports"
)

type logNotifier struct{}

func NewLogNotifier() ports.Notifier {
	return &logNotifier{}
}

func (notifier *logNotifier) NotifyLowStock(product *domain.Product) {
	log.Printf(
		`ALERT: LOW STOCK
		Product Id: %s
		Product Name: %s
		Available Quantity: %d
		Please restock soon to avoid running out of stock.`,
		product.Id, product.Name, product.Quantity)
}
