package main

import (
	"context"
	"io"
	"net"
	"os"
	"time"

	"github.com/brokeyourbike/nickroservices/data"
	"github.com/brokeyourbike/nickroservices/protos"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Currency struct {
	log           hclog.Logger
	rates         *data.ExchangeRates
	subscriptions map[protos.Currency_SubscriberatesServer][]*protos.RateRequest
	protos.UnimplementedCurrencyServer
}

func NewCurrency(l hclog.Logger, r *data.ExchangeRates) *Currency {
	c := &Currency{log: l, rates: r, subscriptions: make(map[protos.Currency_SubscriberatesServer][]*protos.RateRequest)}
	go c.handleUpdates()
	return c
}

func (c *Currency) handleUpdates() {
	ru := c.rates.MonitorRates(5 * time.Second)
	for range ru {
		c.log.Info("Got updated rates")

		// loop over clients
		for k, v := range c.subscriptions {
			// loop over rates
			for _, rr := range v {
				r, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
				if err != nil {
					c.log.Error("Unable to get update rate", "base", rr.GetBase(), "dest", rr.GetDestination(), "error", err)
				}

				err = k.Send(&protos.RateResponse{Base: rr.GetBase(), Destination: rr.GetDestination(), Rate: r})
				if err != nil {
					c.log.Error("Unable to send updated rate", "base", rr.GetBase(), "dest", rr.GetDestination(), "error", err)
				}
			}
		}
	}
}

func (c *Currency) GetRate(ctx context.Context, in *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", in.GetBase(), "destination", in.GetDestination())

	if in.GetBase() == in.GetDestination() {
		s := status.Newf(codes.InvalidArgument, "Base %s cannot be he same as the destination %s", in.GetBase(), in.GetDestination())

		s, err := s.WithDetails(in)
		if err != nil {
			return nil, err
		}

		return nil, s.Err()
	}

	rate, err := c.rates.GetRate(in.GetBase().String(), in.GetDestination().String())
	if err != nil {
		return nil, err
	}

	return &protos.RateResponse{Base: in.GetBase(), Destination: in.GetDestination(), Rate: rate}, nil
}

func (c *Currency) Subscriberates(src protos.Currency_SubscriberatesServer) error {
	for {
		rr, err := src.Recv()
		if err == io.EOF {
			c.log.Info("Client has closed connection")
			break
		}
		if err != nil {
			c.log.Error("Unable to read from client", "error", err)
			break
		}

		c.log.Info("Handle client request", "request", rr)

		rrs, ok := c.subscriptions[src]
		if !ok {
			rrs = []*protos.RateRequest{}
		}

		rrs = append(rrs, rr)
		c.subscriptions[src] = rrs
	}

	return nil
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
