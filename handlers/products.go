package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/brokeyourbike/nickroservices/protos"
	"github.com/hashicorp/go-hclog"
)

type Products struct {
	log hclog.Logger
	cc  protos.CurrencyClient
}

func NewProducts(log hclog.Logger, cc protos.CurrencyClient) *Products {
	return &Products{log: log, cc: cc}
}

func (p *Products) GetProduct(w http.ResponseWriter, h *http.Request) {
	p.log.Info("GetProduct")

	rr := &protos.RateRequest{
		Base:        protos.Currencies_USD,
		Destination: protos.Currencies_EUR,
	}

	resp, err := p.cc.GetRate(context.Background(), rr)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Rate: %f", resp.GetRate())
}
