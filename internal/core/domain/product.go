package domain

import (
	"errors"
	"time"

	"github.com/amangirdhar210/inventory-manager/config"
	"github.com/google/uuid"
)

type Product struct {
	Id        string
	Name      string
	Price     float64
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (product *Product) Validate() error {
	if product.Name == "" {
		return errors.New("product name cannot be empty")
	} else if !isGreaterThanZero(product.Price) {
		return errors.New("product price must be greater than zero")
	} else if product.Quantity < 0 {
		return errors.New("product quantity cannot be negative")
	} else {
		return nil
	}
}

func CreateNewProduct(name string, price float64, quantity int) (*Product, error) {
	product := &Product{
		Id:        uuid.New().String(),
		Name:      name,
		Price:     price,
		Quantity:  quantity,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := product.Validate(); err != nil {
		return nil, err
	}
	return product, nil
}

func isGreaterThanZero[T int | float64](AmountOrQuantity T) bool {
	return AmountOrQuantity > 0
}

func (product *Product) SellUnits(qtyToSell int) error {
	if !isGreaterThanZero(qtyToSell) {
		return errors.New("the quantity to be sold must be greater than zero")
	}

	if product.Quantity < qtyToSell {
		return errors.New("insufficient stock")
	}

	product.Quantity -= qtyToSell

	return nil
}

func (product *Product) Restock(qtyToAdd int) error {
	if !isGreaterThanZero(qtyToAdd) {
		return errors.New("restock amount must be positive")
	}
	product.Quantity += qtyToAdd
	product.UpdatedAt = time.Now()
	return nil
}

func (product *Product) UpdateProductPrice(newPrice float64) error {
	if !isGreaterThanZero(newPrice) {
		return errors.New("price must be greater than zero")
	}
	product.Price = newPrice
	product.UpdatedAt = time.Now()
	return nil
}

func (product *Product) IsLowOnStock() bool {
	return product.Quantity < config.ThresholdAlertQty
}
