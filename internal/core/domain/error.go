package domain

import "errors"

var (
	ErrProductNotFound    = errors.New("product not found")
	ErrProductInvalid     = errors.New("product data is invalid")
	ErrInsufficientStock  = errors.New("insufficient stock")
	ErrRepository         = errors.New("repository error")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrTokenInvalid       = errors.New("token is invalid")
	ErrTokenGeneration    = errors.New("something went wrong while generating token")
)
