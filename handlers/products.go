// Package classification of Product API
//
// Documentation for Product API
// Schemes : http
// BasePath :/
// Version :1.0.0
//
// Consumes:
// - application/json
// Produces:
// - application/json
// swagger:meta

package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/specter25/microservices-in-go/products-api/data"
)

//list of products returns in the response
// swagger:response productResponse
type productsResponse struct {
	//All products in the system
	// in:body
	Body []data.Product
}

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

type KeyProduct struct{}

func (p Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
