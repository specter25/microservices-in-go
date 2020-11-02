package handlers

import (
	"fmt"
	"net/http"

	"github.com/specter25/microservices-in-go/products-api/data"
)

// swagger:route POST /products products createProduct
// Create a new product
//
// responses:
//	200: productResponse
//  422: errorValidation
//  501: errorResponse

// Create handles POST requests to add new products

//AddProduct adds the product received as JSON
func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Post Product")

	prod := r.Context().Value(KeyProduct{}).(data.Product)

	err := prod.Validate()
	if err != nil {
		p.l.Println("[ERROR] Json validation error", err)
		http.Error(rw, fmt.Sprintf("ERROR] Json validation error %s", err), http.StatusBadRequest)
		return
	}

	data.AddProduct(&prod)
	p.l.Printf("Prod %#v", prod)
}
