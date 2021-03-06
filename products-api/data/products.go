package data

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
	protos "github.com/specter25/microservices-in-go/currency-api/protos/currency"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrProductNotFound is an error raised when a product can not be found in the database
var ErrProductNotFound = fmt.Errorf("Product not found")

//Product this is a mock product
//swagger:model
type Product struct {
	// the id of the user
	//
	// required: true
	// min: 1
	ID int `json:"id"`

	// the name for this poduct
	//
	// required: true
	// max length: 255

	Name string `json:"name" validate:"required"`
	// the description for this poduct
	//
	// required: false
	// max length: 10000
	Description string `json:"description"`
	// the price for the product
	//
	// required: true
	// min: 0.01
	Price float64 `json:"price" validate:"required,gt=0"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU       string `json:"sku" validate:"sku"`
	CreatedOn string `json:"-"` // we don;t want these fields to be visible in the output
	UpdatedOn string `json:"-"`
	DeletedOn string `json:"-"`
}

// //Validate validates the jspn data that we receive in the requests
// func (p *Product) Validate() error {
// 	validate := validator.New()
// 	validate.RegisterValidation("sku", validateSKU)
// 	return validate.Struct(p)
// }

type ProductsDB struct {
	currency protos.CurrencyClient
	log      hclog.Logger
	rates    map[string]float64
	client   protos.Currency_SubscribeRatesClient
}

func NewProductDB(c protos.CurrencyClient, l hclog.Logger) *ProductsDB {
	pb := &ProductsDB{c, l, make(map[string]float64), nil}
	go pb.handleUpdates()
	return pb
}

func (p *ProductsDB) handleUpdates() {
	sub, err := p.currency.SubscribeRates(context.Background())
	if err != nil {
		p.log.Error("Unable to subscribe to rates", "error", err)

	}

	p.client = sub

	for {
		rr, err := sub.Recv()
		p.log.Info("Received updated rate from srever", "dest", rr.GetDestination())
		if err != nil {
			p.log.Error("Error receiving message", "error", err)
			return
		}
		p.rates[rr.Destination.String()] = rr.Rate

	}
}

// Products defines a slice of Product
type Products []*Product

// AddProduct adds a new product to the database
func (pr *ProductsDB) AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

// UpdateProduct replaces a product in the database with the given
// item.
// If a product with the given id does not exist in the database
// this function returns a ProductNotFound error
func (pr *ProductsDB) UpdateProduct(id int, p *Product) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}

	productList[pos] = p
	return nil
}
func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, ErrProductNotFound
}

func getNextID() int {
	lp := productList[len(productList)-1]
	value := lp.ID + 1
	return value
}

// GetProducts returns all products from the database
func (p *ProductsDB) GetProducts(currency string) (Products, error) {
	if currency == "" {
		return productList, nil
	}
	// get exchange rate
	rate, err := p.getRate(currency)
	if err != nil {
		p.log.Error("Unable to get the rate ", "currency", currency, "error", err)
		return nil, err
	}

	pr := Products{}
	for _, p := range productList {
		np := *p
		np.Price = np.Price * rate

		pr = append(pr, &np)

	}
	return pr, nil

}

// DeleteProduct deletes a product from the database
func DeleteProduct(id int) error {
	i := findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}

	productList = append(productList[:i], productList[i+1])

	return nil
}

// findIndex finds the index of a product in the database
// returns -1 when no product can be found
func findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

// GetProductByID returns a single product which matches the id from the
// database.
// If a product is not found this function returns a ProductNotFound error
func (p *ProductsDB) GetProductByID(id int, currency string) (*Product, error) {
	i := findIndexByProductID(id)
	if id == -1 {
		return nil, ErrProductNotFound
	}
	if currency == "" {
		return productList[i], nil
	}
	// get exchange rate
	rate, err := p.getRate(currency)
	if err != nil {
		p.log.Error("Unable to get the rate ", "currency", currency, "error", err)
		return nil, err
	}

	np := *productList[i]
	np.Price = np.Price * rate
	return &np, nil
}

func (p *ProductsDB) getRate(destination string) (float64, error) {

	if r, ok := p.rates[destination]; ok {
		return r, nil
	}

	rr := &protos.RateRequest{
		Base:        protos.Currencies(protos.Currencies_value["EUR"]),
		Destination: protos.Currencies(protos.Currencies_value[destination]),
	}

	//Get Initial Rate
	resp, err := p.currency.GetRate(context.Background(), rr)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			md := s.Details()[0].(*protos.RateRequest)
			if s.Code() == codes.InvalidArgument {
				return -1, fmt.Errorf("Unable to get rate from currency server , destination and base currencies cannot be the same , base %s , dest %s", md.Base.String(), md.Destination.String())
			}
			return -1, fmt.Errorf("Unable to get rate from currency server , base %s , dest %s", md.Base.String(), md.Destination.String())
		}
		return -1, err
	}
	p.rates[destination] = resp.Rate

	//Subscripbe for updates
	p.client.Send(rr)

	return resp.Rate, err

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
