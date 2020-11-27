package handlers

import (
	"fmt"
	"log"

	protos "github.com/nicholasjackson/building-microservices-youtube/currency/protos/currency"
	"github.com/specter25/microservices-in-go/products-api/data"
)

// Products handler for getting and updating products
type Products struct {
	l  *log.Logger
	v  *data.Validation
	cc protos.CurrencyClient
}

// NewProducts returns a new products handler with the given logger
func NewProducts(l *log.Logger, v *data.Validation, cc protos.CurrencyClient) *Products {
	return &Products{l, v, cc}
}

// KeyProduct is a key used for the Product object in the context
type KeyProduct struct{}

// ErrInvalidProductPath is an error message when the product path is not valid
var ErrInvalidProductPath = fmt.Errorf("Invalid Path, path should be /products/[id]")

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}
