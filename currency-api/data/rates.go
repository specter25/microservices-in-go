package data

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/go-hclog"
)

//ExchangeRates something
type ExchangeRates struct {
	log  hclog.Logger
	rate map[string]float64
}

//NewRates something
func NewRates(l hclog.Logger) (*ExchangeRates, error) {
	er := &ExchangeRates{log: l, rate: map[string]float64{}}

	er.getRates()

	return er, nil
}

// GetRate something
func (e *ExchangeRates) GetRate(base string, dest string) (float64, error) {
	br, ok := e.rate[base]
	if !ok {
		return 0, fmt.Errorf("rate not found for currency %s", base)
	}
	dr, ok2 := e.rate[dest]
	if !ok2 {
		return 0, fmt.Errorf("Rate not found for currency %s", dest)
	}
	return dr / br, nil
}

func (e *ExchangeRates) getRates() error {
	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected error code 200 got %d", resp.StatusCode)

	}
	defer resp.Body.Close()

	nd := &Cubes{}

	xml.NewDecoder(resp.Body).Decode(&nd)

	for _, c := range nd.CubeData {
		r, err := strconv.ParseFloat(c.Rate, 64)
		if err != nil {
			return err
		}
		e.rate[c.Currency] = r

	}

	e.rate["EUR"] = 1

	return nil
}

//Cubes something
type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

//Cube something
type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}
