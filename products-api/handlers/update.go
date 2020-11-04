package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/specter25/microservices-in-go/products-api/data"
)

// swagger:route PUT /products products updateProduct
// Update a products details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  422: errorValidation

// UpdateProducts handles PUT requests to update products
func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {

	//http req body is an ioreader
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert to integer", http.StatusBadRequest)
	}
	p.l.Println("Handle Post Product", id)

	rw.Header().Add("Content-Type", "application/json")

	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	p.l.Println("[DEBUG] updating record id", prod.ID)

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
	if err != nil {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	p.l.Printf("Prod %#v", prod)
}
