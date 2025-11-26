package dto

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AddProductRequest struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type UpdateQuantityRequest struct {
	Quantity int `json:"quantity"`
}

type UpdatePriceRequest struct {
	Price float64 `json:"price"`
}
