package main

import (
	"context"
	"net"
	"os"

	"github.com/brokeyourbike/nickroservices/protos"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Currency struct {
	log hclog.Logger
	protos.UnimplementedCurrencyServer
}

func NewCurrency(l hclog.Logger) *Currency {
	return &Currency{log: l}
}

func (c *Currency) GetRate(ctx context.Context, in *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", in.GetBase(), "destination", in.GetDestination())

	return &protos.RateResponse{Rate: 0.5}, nil
}

func main() {
	log := hclog.Default()

	l, err := net.Listen("tcp", "127.0.0.1:9092")
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}

	gs := grpc.NewServer()
	cs := NewCurrency(log)

	protos.RegisterCurrencyServer(gs, cs)

	reflection.Register(gs)
	gs.Serve(l)
}
