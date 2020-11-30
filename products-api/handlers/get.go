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
	p.l.Debug("Get all records")
	rw.Header().Add("Content-Type", "application/json")

	prods, err := p.productDB.GetProducts("")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}
	//This is one way to convert it into json
	// d, err := json.Marshal(prods)
	// Now the second way
	err = data.ToJSON(prods, rw)

	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
		p.l.Error("Unable to serialize the products", "error", err)

	}
	//we needed this in the forst approach not now
	// rw.Write(d)
}

// swagger:route GET /products/{id} products listSingleProduct
// Return a list of products from the database
// responses:
//	200: productResponse
//	404: errorResponse

// ListSingle handles GET requests
func (p *Products) ListSingle(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	id := getProductID(r)

	p.l.Debug("Get record id", id)

	prod, err := p.productDB.GetProductByID(id, "")

	switch err {
	case nil:

	case data.ErrProductNotFound:
		p.l.Error("fetching product", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		p.l.Error("fetching product", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJSON(prod, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Error("serializing product", err)
	}
}
