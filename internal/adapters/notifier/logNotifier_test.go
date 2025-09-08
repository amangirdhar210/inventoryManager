package notifier

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/amangirdhar210/inventory-manager/internal/core/domain"
)

func TestLogNotifier_NotifyLowStock(t *testing.T) {
	notifier := NewLogNotifier()
	product := &domain.Product{
		Id:       "prod-abc-123",
		Name:     "Gaming Mouse",
		Quantity: 5,
	}

	var buf bytes.Buffer
	originalOutput := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(originalOutput)

	notifier.NotifyLowStock(product)

	output := buf.String()

	expectedSubstrings := []string{
		"ALERT: LOW STOCK",
		"Product Id: prod-abc-123",
		"Product Name: Gaming Mouse",
		"Available Quantity: 5",
	}

	for _, sub := range expectedSubstrings {
		if !strings.Contains(output, sub) {
			t.Errorf("log output did not contain expected substring %q. Full output: %q", sub, output)
		}
	}
}
