package main

import (
	"context"
	"net"
	"os"

	"github.com/brokeyourbike/nickroservices/data"
	"github.com/brokeyourbike/nickroservices/protos"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Currency struct {
	log   hclog.Logger
	rates *data.ExchangeRates
	protos.UnimplementedCurrencyServer
}

func NewCurrency(l hclog.Logger, r *data.ExchangeRates) *Currency {
	return &Currency{log: l, rates: r}
}

func (c *Currency) GetRate(ctx context.Context, in *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", in.GetBase(), "destination", in.GetDestination())

	rate, err := c.rates.GetRate(in.GetBase().String(), in.GetDestination().String())
	if err != nil {
		return nil, err
	}

	return &protos.RateResponse{Rate: rate}, nil
}

func main() {
	log := hclog.Default()

	rates, err := data.NewExchangeRates(log)
	if err != nil {
		log.Error("Unable to instantiate rates", "error", err)
		os.Exit(1)
	}
	cs := NewCurrency(log, rates)

	l, err := net.Listen("tcp", "127.0.0.1:9092")
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}

	gs := grpc.NewServer()
	protos.RegisterCurrencyServer(gs, cs)
	reflection.Register(gs)
	gs.Serve(l)
}
