package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/specter25/microservices-in-go/products-api/data"
)

//UpdateProducts updates the products based on data received
func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {

	//http req body is an ioreader

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert to integer", http.StatusBadRequest)
	}
	p.l.Println("Handle Post Product", id)

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	err = prod.Validate()
	if err != nil {
		p.l.Println("[ERROR] Json validation error", err)
		http.Error(rw, fmt.Sprintf("ERROR] Json validation error %s", err), http.StatusBadRequest)
		return
	}

	err = data.UpdateProduct(id, &prod)
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
