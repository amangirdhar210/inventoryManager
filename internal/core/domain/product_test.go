package domain

import (
	"testing"
)

const ThresholdAlertQty = 10

func TestNewProduct(t *testing.T) {

	tests := []struct {
		name        string
		productName string
		price       float64
		quantity    int
		expectErr   bool
	}{
		{"should create product successfully", "Macbook Pro", 1500.00, 20, false},
		{"should fail with empty name", "", 1500.00, 20, true},
		{"should fail with zero price", "Macbook Pro", 0, 20, true},
		{"should fail with negative price", "Macbook Pro", -1, 20, true},
		{"should fail with negative quantity", "Macbook Pro", 1500.00, -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := NewProduct(tt.productName, tt.price, tt.quantity)

			if (err != nil) != tt.expectErr {
				t.Errorf("NewProduct() error = %v, expectErr %v", err, tt.expectErr)
				return
			}

			if !tt.expectErr {
				if product == nil {
					t.Fatal("NewProduct() returned nil product on success")
				}
				if product.Name != tt.productName || product.Price != tt.price || product.Quantity != tt.quantity {
					t.Errorf("NewProduct() got = %+v, want name=%s, price=%f, quantity=%d", product, tt.productName, tt.price, tt.quantity)
				}
				if product.Id == "" {
					t.Error("NewProduct() did not assign an ID")
				}
			}
		})
	}
}

func TestProduct_IsLowOnStock(t *testing.T) {

	tests := []struct {
		name          string
		quantity      int
		expectedIsLow bool
	}{
		{"should be true when quantity is below threshold", ThresholdAlertQty - 1, true},
		{"should be true when quantity is zero", 0, true},
		{"should be false when quantity is at threshold", ThresholdAlertQty, false},
		{"should be false when quantity is above threshold", ThresholdAlertQty + 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := &Product{Quantity: tt.quantity}

			if got := p.IsLowOnStock(); got != tt.expectedIsLow {
				t.Errorf("IsLowOnStock() = %v, want %v", got, tt.expectedIsLow)
			}
		})
	}
}
