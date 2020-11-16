package main

import (
	"net"
	"os"

	"github.com/hashicorp/go-hclog"
	protos "github.com/specter25/microservices-in-go/currency-api/protos/currency"
	"github.com/specter25/microservices-in-go/currency-api/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()

	gs := grpc.NewServer()
	cs := server.NewCurrency(log)

	protos.RegisterCurrencyServer(gs, cs)

	reflection.Register(gs)

	l, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}

	gs.Serve(l)

}
