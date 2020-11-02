package handlers

import (
	"net/http"

	"github.com/specter25/microservices-in-go/products-api/data"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
// 		200: productResponses

//GetProducts returns the products from the array
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
