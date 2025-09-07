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
		`ALERT: LOW STOCK \n
		Product Id: %s \n
		Product Name: %s \n
		Available Quantity: %d \n
		Please restock soon to avoid running out of stock.\n`,
		product.Id, product.Name, product.Quantity)
}
