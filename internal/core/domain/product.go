package domain

import (
	"github.com/google/uuid"
)

type Product struct {
	Id       string  `dynamodbav:"id"`
	Name     string  `dynamodbav:"name"`
	Price    float64 `dynamodbav:"price"`
	Quantity int     `dynamodbav:"quantity"`
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
