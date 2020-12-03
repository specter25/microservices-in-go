package server

import (
	"context"
	"io"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/go-hclog"
	data "github.com/specter25/microservices-in-go/currency-api/data"
	protos "github.com/specter25/microservices-in-go/currency-api/protos/currency"
)

type Currency struct {
	rates         *data.ExchangeRates
	log           hclog.Logger
	subscriptions map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest
}

func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {
	c := &Currency{r, l, make(map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest)}
	go c.HandleUpdates()

	return c
}

func (c *Currency) HandleUpdates() {
	ru := c.rates.MonitorRates(5 * time.Second)
	for range ru {
		c.log.Info("Got updated rates")
		//loop over subscribed clients
		for k, v := range c.subscriptions {
			//loop over subscribed rates
			for _, rr := range v {
				r, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
				if err != nil {
					c.log.Error("Unable to get updated rates", "base", rr.GetBase().String(), "destination", rr.GetDestination().String())
				}
				err = k.Send(&protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: r})
				if err != nil {
					c.log.Error("Unable to send updated rates", "base", rr.GetBase().String(), "destination", rr.GetDestination().String())
				}
			}

		}
	}

}

func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle Get Rate", "base", rr.GetBase(), "destination", rr.GetDestination())

	if rr.Base == rr.Destination {
		err := status.Newf(
			codes.InvalidArgument,
			"Base currency %s cannot be the same as the destination currency %s",
			rr.Base.String(),
			rr.Destination.String())
		//adding metadata to the error object
		err, wde := err.WithDetails(rr)
		if wde != nil {
			return nil, wde
		}
		return nil, err.Err()
	}

	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}
	return &protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: rate}, nil
}

func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {

	for {
		rr, err := src.Recv()
		if err == io.EOF {
			c.log.Error("Client has closed connections")
			break
		}
		if err != nil {
			c.log.Error("unable to read form client ", "error", err)
			return err
		}
		c.log.Info("handle client request", "request_base", rr.GetBase(), "request_destination", rr.GetDestination())
		rrs, ok := c.subscriptions[src]
		if !ok {
			rrs = []*protos.RateRequest{}
		}
		rrs = append(rrs, rr)
		c.subscriptions[src] = rrs
	}

	return nil

}
