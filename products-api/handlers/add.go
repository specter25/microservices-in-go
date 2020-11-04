package handlers

import (
	"net/http"

	"github.com/specter25/microservices-in-go/data"
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

	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	data.AddProduct(prod)
	p.l.Printf("Prod %#v", prod)
}
