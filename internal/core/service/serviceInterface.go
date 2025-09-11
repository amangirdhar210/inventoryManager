package service

import "github.com/amangirdhar210/inventory-manager/internal/core/domain"

type InventoryService interface {
	AddProduct(name string, price float64, quantity int) (*domain.Product, error)
	GetProduct(id string) (*domain.Product, error)
	SellProductUnits(id string, quantity int) error
	RestockProduct(id string, quantity int) error
	UpdateProductPrice(id string, newPrice float64) error
	GetAllProducts() ([]domain.Product, error)
	DeleteProduct(id string) error
	GetInventoryValue() (float64, error)
}

type AuthService interface {
	Login(email, password string) (string, error)
}
