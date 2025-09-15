package domain

import (
	"github.com/google/uuid"
)

type Product struct {
	Id       string
	Name     string
	Price    float64
	Quantity int
}

func NewProduct(name string, price float64, quantity int) *Product {
	product := &Product{
		Id:       uuid.New().String(),
		Name:     name,
		Price:    price,
		Quantity: quantity,
	}
	return product
}
