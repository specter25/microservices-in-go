package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

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
	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}
	//handle an Update
	if r.Method == http.MethodPut {
		// now expect the id in the URI
		// now we ahve to use regex to have the :id of te url
		re := regexp.MustCompile(`/([0-9]+)`)
		g := re.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			http.Error(rw, "Invalid URL", http.StatusBadRequest)
			return
		}
		if len(g[0]) != 2 {
			http.Error(rw, "Invalid URL 2", http.StatusBadRequest)
			return
		}
		idString := g[0][1]
		id, err := strconv.Atoi(idString)

		if err != nil {
			http.Error(rw, "Invalid string", http.StatusBadRequest)
			return
		}

		p.l.Println("got id", id)
		p.updateProducts(id, rw, r)

	}

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

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Post Product")
	//http req body is an ioreader
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}
	data.AddProduct(prod)
	p.l.Printf("Prod %#v", prod)
}

func (p *Products) updateProducts(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Post Product")
	//http req body is an ioreader
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}
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
