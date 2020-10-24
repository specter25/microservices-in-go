package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/specter25/microservices-in-go/products-api/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	// make a get request to srev .http and return this product list
	// we need ro look at a package known as encoding.json
	// this is convert the product struct into 2 json
	//there are 2 ways to go about it

	lp := data.GetProducts()
	//This is one way to convert it into json
	// d, err := json.Marshal(lp)
	// Now the second way
	err := lp.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
	//we needed this in the forst approach not now
	// rw.Write(d)
}

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
