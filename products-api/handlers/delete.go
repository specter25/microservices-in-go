package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/specter25/microservices-in-go/data"
)

// swagger:route DELETE /products/{id} products deleteProducts
// Returns a list of products
// responses:
// 		200: noContent

//DeleteProduct deletes a product from the database
func (p *Products) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	p.l.Println("[DEBUG] deleting record id", id)

	err := data.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		p.l.Println("[ERROR] deleting record id does not exist")

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	if err != nil {
		p.l.Println("[ERROR] deleting record", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)

}
