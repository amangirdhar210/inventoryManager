package domain

import (
	"errors"

	"github.com/amangirdhar210/inventory-manager/config"
	"github.com/google/uuid"
)

type Product struct {
	Id       string
	Name     string
	Price    float64
	Quantity int
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

func NewProduct(name string, price float64, quantity int) (*Product, error) {
	product := &Product{
		Id:       uuid.New().String(),
		Name:     name,
		Price:    price,
		Quantity: quantity,
	}

	if err := product.Validate(); err != nil {
		return nil, err
	}
	return product, nil
}

func isGreaterThanZero[T int | float64](AmountOrQuantity T) bool {
	return AmountOrQuantity > 0
}

func (product *Product) IsLowOnStock() bool {
	return product.Quantity < config.ThresholdAlertQty
}
