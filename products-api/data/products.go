package data

import (
	"encoding/json"
	"io"
	"time"
)

//Product this is a mock product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"` // we don;t want these fields to be visible in the output
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

type Products []*Product

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

func getNextID() int {
	lp := productList[len(productList)-1]
	value := lp.ID + 1
	return value
}

func GetProducts() Products {
	return productList
}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		SKU:         "abc323",
		Price:       2.45,
		UpdatedOn:   time.Now().UTC().String(),
		CreatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Expresso",
		Description: "Short and Strong Coffee without milk",
		SKU:         "fjd34",
		Price:       1.99,
		UpdatedOn:   time.Now().UTC().String(),
		CreatedOn:   time.Now().UTC().String(),
	},
}
