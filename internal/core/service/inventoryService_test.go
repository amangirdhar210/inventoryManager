package service

import (
	"errors"
	"testing"

	"github.com/amangirdhar210/inventory-manager/internal/core/domain"
)

var (
	ErrRepoFailed = errors.New("repository failed")
)

type mockProductRepository struct {
	products    map[string]*domain.Product
	shouldError bool
}

func newMockProductRepository() *mockProductRepository {
	return &mockProductRepository{
		products: make(map[string]*domain.Product),
	}
}

func (m *mockProductRepository) Save(product *domain.Product) error {
	if m.shouldError {
		return ErrRepoFailed
	}
	m.products[product.Id] = product
	return nil
}

func (m *mockProductRepository) FindById(id string) (*domain.Product, error) {
	if m.shouldError {
		return nil, ErrRepoFailed
	}
	product, ok := m.products[id]
	if !ok {
		return nil, errors.New("product not found")
	}
	clone := *product
	return &clone, nil
}

func (m *mockProductRepository) Update(product *domain.Product) error {
	if m.shouldError {
		return ErrRepoFailed
	}
	if _, ok := m.products[product.Id]; !ok {
		return errors.New("product not found for update")
	}
	m.products[product.Id] = product
	return nil
}

func (m *mockProductRepository) ListAll() ([]domain.Product, error) {
	if m.shouldError {
		return nil, ErrRepoFailed
	}
	var productList []domain.Product
	for _, p := range m.products {
		productList = append(productList, *p)
	}
	return productList, nil
}

func (m *mockProductRepository) DeleteById(id string) error {
	if m.shouldError {
		return ErrRepoFailed
	}
	if _, ok := m.products[id]; !ok {
		return errors.New("product not found for deletion")
	}
	delete(m.products, id)
	return nil
}

type mockNotifier struct {
	notifiedProduct *domain.Product
	wasCalled       bool
}

func (m *mockNotifier) NotifyLowStock(product *domain.Product) {
	m.wasCalled = true
	m.notifiedProduct = product
}

func TestInventoryService_AddProduct(t *testing.T) {
	tests := []struct {
		name        string
		productName string
		price       float64
		quantity    int
		repoShould  bool
		expectErr   bool
	}{
		{"success", "Laptop", 1200.00, 10, false, false},
		{"fail_invalid_name", "", 1200.00, 10, false, true},
		{"fail_invalid_price", "Laptop", -1, 10, false, true},
		{"fail_repo_save", "Laptop", 1200.00, 10, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockProductRepository()
			repo.shouldError = tt.repoShould
			service := NewInventoryService(repo, &mockNotifier{})

			product, err := service.AddProduct(tt.productName, tt.price, tt.quantity)

			if (err != nil) != tt.expectErr {
				t.Errorf("AddProduct() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && (product == nil || len(repo.products) != 1) {
				t.Errorf("AddProduct() failed to create or save the product")
			}
		})
	}
}

func TestInventoryService_GetProduct(t *testing.T) {
	repo := newMockProductRepository()
	p, _ := domain.NewProduct("Test Book", 25.50, 50)
	repo.Save(p)

	tests := []struct {
		name      string
		productID string
		expectErr bool
	}{
		{"success", p.Id, false},
		{"fail_not_found", "non-existent-id", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewInventoryService(repo, &mockNotifier{})
			product, err := service.GetProduct(tt.productID)

			if (err != nil) != tt.expectErr {
				t.Errorf("GetProduct() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && product.Id != tt.productID {
				t.Errorf("GetProduct() got = %v, want %v", product.Id, tt.productID)
			}
		})
	}
}

func TestInventoryService_SellProductUnits(t *testing.T) {
	p, _ := domain.NewProduct("Monitor", 300, 20)

	tests := []struct {
		name           string
		initialProduct *domain.Product
		sellQuantity   int
		repoShould     bool
		expectErr      bool
		notifierCalled bool
		finalQuantity  int
	}{
		{"success", p, 5, false, false, false, 15},
		{"success_low_stock_notification", p, 11, false, false, true, 9},
		{"fail_insufficient_stock", p, 25, false, true, false, 20},
		{"fail_product_not_found", p, 5, false, true, false, 0},
		{"fail_repo_update", p, 5, true, true, false, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockProductRepository()
			if tt.name != "fail_product_not_found" {
				clone := *tt.initialProduct
				repo.Save(&clone)
			}
			repo.shouldError = tt.repoShould
			notifier := &mockNotifier{}
			service := NewInventoryService(repo, notifier)

			productID := tt.initialProduct.Id
			if tt.name == "fail_product_not_found" {
				productID = "wrong-id"
			}

			err := service.SellProductUnits(productID, tt.sellQuantity)

			if (err != nil) != tt.expectErr {
				t.Errorf("SellProductUnits() error = %v, expectErr %v", err, tt.expectErr)
			}

			if !tt.expectErr {
				updatedProduct, _ := repo.FindById(productID)
				if updatedProduct.Quantity != tt.finalQuantity {
					t.Errorf("SellProductUnits() final quantity = %d, want %d", updatedProduct.Quantity, tt.finalQuantity)
				}
			}

			if notifier.wasCalled != tt.notifierCalled {
				t.Errorf("SellProductUnits() notifier called = %v, want %v", notifier.wasCalled, tt.notifierCalled)
			}
		})
	}
}

func TestInventoryService_RestockProduct(t *testing.T) {
	p, _ := domain.NewProduct("Keyboard", 75, 10)

	tests := []struct {
		name          string
		restockQty    int
		repoShould    bool
		expectErr     bool
		finalQuantity int
	}{
		{"success", 20, false, false, 30},
		{"fail_invalid_quantity", -5, false, true, 10},
		{"fail_repo_update", 20, true, true, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockProductRepository()
			clone := *p
			repo.Save(&clone)
			repo.shouldError = tt.repoShould
			service := NewInventoryService(repo, &mockNotifier{})

			err := service.RestockProduct(p.Id, tt.restockQty)

			if (err != nil) != tt.expectErr {
				t.Errorf("RestockProduct() error = %v, expectErr %v", err, tt.expectErr)
			}

			if !tt.expectErr {
				updatedProduct, _ := repo.FindById(p.Id)
				if updatedProduct.Quantity != tt.finalQuantity {
					t.Errorf("RestockProduct() final quantity = %d, want %d", updatedProduct.Quantity, tt.finalQuantity)
				}
			}
		})
	}
}

func TestInventoryService_GetAllProducts(t *testing.T) {
	p1, _ := domain.NewProduct("Product A", 10, 1)
	p2, _ := domain.NewProduct("Product B", 20, 2)

	tests := []struct {
		name       string
		setupRepo  func() *mockProductRepository
		wantCount  int
		expectErr  bool
		wantResult []domain.Product
	}{
		{
			"success_with_products",
			func() *mockProductRepository {
				repo := newMockProductRepository()
				repo.Save(p1)
				repo.Save(p2)
				return repo
			},
			2, false, []domain.Product{*p1, *p2},
		},
		{
			"success_no_products",
			func() *mockProductRepository {
				return newMockProductRepository()
			},
			0, false, []domain.Product{},
		},
		{
			"fail_repo_error",
			func() *mockProductRepository {
				repo := newMockProductRepository()
				repo.shouldError = true
				return repo
			},
			0, true, nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.setupRepo()
			service := NewInventoryService(repo, &mockNotifier{})
			products, err := service.GetAllProducts()

			if (err != nil) != tt.expectErr {
				t.Errorf("GetAllProducts() error = %v, expectErr %v", err, tt.expectErr)
			}

			if !tt.expectErr && len(products) != tt.wantCount {
				t.Errorf("GetAllProducts() count = %d, want %d", len(products), tt.wantCount)
			}
		})
	}
}

func TestInventoryService_DeleteProduct(t *testing.T) {
	p, _ := domain.NewProduct("ToDelete", 1, 1)

	tests := []struct {
		name      string
		productID string
		setupRepo func() *mockProductRepository
		expectErr bool
	}{
		{
			"success", p.Id,
			func() *mockProductRepository {
				repo := newMockProductRepository()
				repo.Save(p)
				return repo
			},
			false,
		},
		{
			"fail_not_found", "wrong-id",
			func() *mockProductRepository {
				repo := newMockProductRepository()
				repo.Save(p)
				return repo
			},
			true,
		},
		{
			"fail_repo_error", p.Id,
			func() *mockProductRepository {
				repo := newMockProductRepository()
				repo.Save(p)
				repo.shouldError = true
				return repo
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.setupRepo()
			service := NewInventoryService(repo, &mockNotifier{})
			err := service.DeleteProduct(tt.productID)

			if (err != nil) != tt.expectErr {
				t.Errorf("DeleteProduct() error = %v, expectErr %v", err, tt.expectErr)
			}
			if !tt.expectErr {
				if _, ok := repo.products[tt.productID]; ok {
					t.Error("DeleteProduct() failed to remove product from repo")
				}
			}
		})
	}
}

func TestInventoryService_GetInventoryValue(t *testing.T) {
	p1, _ := domain.NewProduct("Valuable", 10.50, 10)
	p2, _ := domain.NewProduct("Cheap", 1.00, 100)

	tests := []struct {
		name      string
		setupRepo func() *mockProductRepository
		wantValue float64
		expectErr bool
	}{
		{
			"success",
			func() *mockProductRepository {
				repo := newMockProductRepository()
				repo.Save(p1)
				repo.Save(p2)
				return repo
			},
			205.00, false,
		},
		{
			"success_empty",
			func() *mockProductRepository {
				return newMockProductRepository()
			},
			0, false,
		},
		{
			"fail_repo_error",
			func() *mockProductRepository {
				repo := newMockProductRepository()
				repo.shouldError = true
				return repo
			},
			0, true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.setupRepo()
			service := NewInventoryService(repo, &mockNotifier{})
			value, err := service.GetInventoryValue()

			if (err != nil) != tt.expectErr {
				t.Errorf("GetInventoryValue() error = %v, expectErr %v", err, tt.expectErr)
			}
			if !tt.expectErr && value != tt.wantValue {
				t.Errorf("GetInventoryValue() got = %f, want %f", value, tt.wantValue)
			}
		})
	}
}
