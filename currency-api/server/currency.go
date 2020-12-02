package server

import (
	"context"
	"io"
	"time"

	"github.com/hashicorp/go-hclog"
	data "github.com/specter25/microservices-in-go/currency-api/data"
	protos "github.com/specter25/microservices-in-go/currency-api/protos/currency"
)

type Currency struct {
	rates *data.ExchangeRates
	log   hclog.Logger
}

func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {
	return &Currency{r, l}
}

func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle Get Rate", "base", rr.GetBase(), "destination", rr.GetDestination())

	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}
	return &protos.RateResponse{Rate: rate}, nil
}

func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {

	go func() {
		for {
			rr, err := src.Recv()
			if err == io.EOF {
				c.log.Error("Client has closed connections")
				break
			}
			if err != nil {
				c.log.Error("unable to read form client ", "error", err)
				break
			}
			c.log.Info("handle client request", "request_base", rr.GetBase(), "request_destination", rr.GetDestination())
		}
	}()

	// to handle the message send functionality
	for {
		err := src.Send(&protos.RateResponse{Rate: 12.1})
		if err != nil {
			return err
		}
		time.Sleep(5 * time.Second)
	}
}
