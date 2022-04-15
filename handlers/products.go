package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/brokeyourbike/nickroservices/protos"
	"github.com/hashicorp/go-hclog"
)

type Products struct {
	log    hclog.Logger
	cc     protos.CurrencyClient
	client protos.Currency_SubscriberatesClient
	rates  map[string]float64
}

func NewProducts(log hclog.Logger, cc protos.CurrencyClient) *Products {
	p := &Products{log: log, cc: cc, client: nil, rates: make(map[string]float64)}
	go p.handleUpdates()
	return p
}

func (p *Products) handleUpdates() {
	sub, err := p.cc.Subscriberates(context.Background())
	if err != nil {
		p.log.Error("Cannot subscribe for rates", "error", err)
	}

	p.client = sub

	for {
		resp, err := sub.Recv()
		if err != nil {
			p.log.Error("Error receiving message", "error", err)
			return
		}

		p.log.Info("Received updated rate", "dest", resp.Destination, "rate", resp.GetRate())
		p.rates[resp.Destination.String()] = resp.GetRate()
	}
}

func (p *Products) GetProduct(w http.ResponseWriter, h *http.Request) {
	p.log.Info("GetProduct")

	rate, err := p.getRateFor(protos.Currencies_USD, protos.Currencies_EUR)
	if err != nil {
		p.log.Error("Unable to get rate", "error", err)
		http.Error(w, "Service unavailable", http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "Rate: %f", rate)
}

func (p *Products) getRateFor(base, dest protos.Currencies) (float64, error) {
	if r, ok := p.rates[dest.String()]; ok {
		return r, nil
	}

	rr := &protos.RateRequest{
		Base:        base,
		Destination: dest,
	}

	// get initial rate
	resp, err := p.cc.GetRate(context.Background(), rr)
	if err != nil {
		return 0, err
	}

	p.rates[resp.Destination.String()] = resp.GetRate()

	err = p.client.Send(rr)

	return resp.GetRate(), err
}
