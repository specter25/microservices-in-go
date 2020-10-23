package handlers

import (
	"log"
	"net/http"

	"github.com/specter25/microservices-in-go/products-api/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	//handle an Update

	rw.WriteHeader(http.StatusMethodNotAllowed)

}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
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
