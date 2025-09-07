package ports

import "github.com/amangirdhar210/inventory-manager/internal/core/domain"

type InventoryService interface {
	AddProduct(name string, price float64, quantity int) (*domain.Product, error)
	GetProduct(id string) (*domain.Product, error)
	SellProductUnits(id string, quantity int) (*domain.Product, error)
	RestockProduct(id string, quantity int) (*domain.Product, error)
	GetAllProducts() ([]domain.Product, error)
	DeleteProduct(id string) error
	GetInventoryValue() (float64, error)
}

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
