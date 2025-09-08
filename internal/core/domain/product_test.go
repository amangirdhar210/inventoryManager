package domain

import (
	"testing"
)

const ThresholdAlertQty = 10

func TestCreateNewProduct(t *testing.T) {

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
			product, err := CreateNewProduct(tt.productName, tt.price, tt.quantity)

			if (err != nil) != tt.expectErr {
				t.Errorf("CreateNewProduct() error = %v, expectErr %v", err, tt.expectErr)
				return
			}

			if !tt.expectErr {
				if product == nil {
					t.Fatal("CreateNewProduct() returned nil product on success")
				}
				if product.Name != tt.productName || product.Price != tt.price || product.Quantity != tt.quantity {
					t.Errorf("CreateNewProduct() got = %+v, want name=%s, price=%f, quantity=%d", product, tt.productName, tt.price, tt.quantity)
				}
				if product.Id == "" {
					t.Error("CreateNewProduct() did not assign an ID")
				}
			}
		})
	}
}

func TestProduct_SellUnits(t *testing.T) {

	tests := []struct {
		name        string
		initialQty  int
		sellQty     int
		expectedQty int
		expectErr   bool
	}{
		{"should sell units successfully", 20, 5, 15, false},
		{"should sell all remaining units", 20, 20, 0, false},
		{"should fail for insufficient stock", 5, 10, 5, true},
		{"should fail for zero quantity", 20, 0, 20, true},
		{"should fail for negative quantity", 20, -5, 20, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Product{Id: "test-id", Name: "Test Product", Price: 100, Quantity: tt.initialQty}

			err := p.SellUnits(tt.sellQty)

			if (err != nil) != tt.expectErr {
				t.Errorf("SellUnits() error = %v, expectErr %v", err, tt.expectErr)
			}
			if p.Quantity != tt.expectedQty {
				t.Errorf("SellUnits() quantity got = %v, want %v", p.Quantity, tt.expectedQty)
			}
		})
	}
}

func TestProduct_Restock(t *testing.T) {

	tests := []struct {
		name        string
		initialQty  int
		addQty      int
		expectedQty int
		expectErr   bool
	}{
		{"should restock successfully", 10, 15, 25, false},
		{"should fail for zero quantity", 10, 0, 10, true},
		{"should fail for negative quantity", 10, -5, 10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := &Product{Id: "test-id", Name: "Test Product", Price: 100, Quantity: tt.initialQty}

			err := p.Restock(tt.addQty)

			if (err != nil) != tt.expectErr {
				t.Errorf("Restock() error = %v, expectErr %v", err, tt.expectErr)
			}

			if p.Quantity != tt.expectedQty {
				t.Errorf("Restock() quantity got = %v, want %v", p.Quantity, tt.expectedQty)
			}
		})
	}
}

func TestProduct_UpdateProductPrice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		initialPrice  float64
		newPrice      float64
		expectedPrice float64
		expectErr     bool
	}{
		{"should update price successfully", 50.0, 75.50, 75.50, false},
		{"should fail for zero price", 50.0, 0, 50.0, true},
		{"should fail for negative price", 50.0, -10.0, 50.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := &Product{Id: "test-id", Name: "Test Product", Price: tt.initialPrice, Quantity: 10}

			err := p.UpdateProductPrice(tt.newPrice)

			if (err != nil) != tt.expectErr {
				t.Errorf("UpdateProductPrice() error = %v, expectErr %v", err, tt.expectErr)
			}
			if p.Price != tt.expectedPrice {
				t.Errorf("UpdateProductPrice() price got = %v, want %v", p.Price, tt.expectedPrice)
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

//implement testing for UpdateProductPrice
