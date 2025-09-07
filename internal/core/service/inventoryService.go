package service

import (
	"fmt"

	"github.com/amangirdhar210/inventory-manager/internal/core/domain"
	"github.com/amangirdhar210/inventory-manager/internal/core/ports"
)

type inventoryService struct {
	repo     ports.ProductRepository
	notifier ports.Notifier
}

func NewInventoryService(repo ports.ProductRepository, notifier ports.Notifier) ports.InventoryService {
	return &inventoryService{
		repo:     repo,
		notifier: notifier,
	}
}

func (invService *inventoryService) AddProduct(name string, price float64, quantity int) (*domain.Product, error) {
	product, err := domain.CreateNewProduct(name, price, quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to create new product : %w", err)
	}

	if err := invService.repo.Save(product); err != nil {
		return nil, fmt.Errorf("failed to save product: %w ", err)
	}

	return product, nil
}

func (invService *inventoryService) GetProduct(id string) (*domain.Product, error) {
	product, err := invService.repo.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product with id %s: %w", id, err)
	}
	return product, nil
}

func (invService *inventoryService) SellProductUnits(id string, quantity int) (*domain.Product, error) {
	product, err := invService.repo.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("could not find the product for sale: %w", err)
	}

	if err := product.SellUnits(quantity); err != nil {
		return nil, fmt.Errorf("failed to sell the product: %w", err)
	}

	if err := invService.repo.Update(product); err != nil {
		return nil, fmt.Errorf("failed to update product stock after sale: %w", err)
	}

	if product.IsLowOnStock() {
		invService.notifier.NotifyLowStock(product)
	}
	return product, nil
}

func (invService *inventoryService) RestockProduct(id string, quantity int) (*domain.Product, error) {
	product, err := invService.repo.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("could not find the product to be restocked: %w", err)
	}

	if err := product.Restock(quantity); err != nil {
		return nil, fmt.Errorf("failed to restock the product: %w", err)
	}

	if err := invService.repo.Update(product); err != nil {
		return nil, fmt.Errorf("failed to update product stock after restock: %w", err)
	}

	return product, nil

}

func (invService *inventoryService) GetAllProducts() ([]domain.Product, error) {
	products, err := invService.repo.ListAll()
	if err != nil {
		return nil, fmt.Errorf("failed to list all products: %w", err)
	}
	return products, nil
}

func (invService *inventoryService) DeleteProduct(id string) error {
	err := invService.repo.DeleteById(id)
	if err != nil {
		return fmt.Errorf("failed to delete product with id %s: %w", id, err)
	}
	return nil
}

func (invService *inventoryService) GetInventoryValue() (float64, error) {
	products, err := invService.repo.ListAll()
	if err != nil {
		return 0, fmt.Errorf("failed to list all products: %w", err)
	}

	var totalValue float64 = 0
	for _, product := range products {
		totalValue += float64(product.Quantity) * product.Price
	}
	return totalValue, nil
}

func (invService *inventoryService) UpdateProductPrice(id string, newPrice float64) error {
	product, err := invService.repo.FindById(id)
	if err != nil {
		return fmt.Errorf("%w: could not find product with id %s", domain.ErrProductNotFound, id)
	}

	err = product.UpdateProductPrice(newPrice)
	if err != nil {
		return fmt.Errorf("failed to update price of the product: %w", err)
	}

	err = invService.repo.Update(product)
	if err != nil {
		return fmt.Errorf("could not save the updated price: %w", err)
	}

	return nil
}
