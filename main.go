package main

import (
	"net"
	"os"

	"github.com/brokeyourbike/nickroservices/protos"
	"github.com/brokeyourbike/nickroservices/server"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()

	l, err := net.Listen("tcp", "127.0.0.1:9092")
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}

	gs := grpc.NewServer()
	cs := server.NewCurrency(log)

	protos.RegisterCurrencyServer(gs, cs)

	reflection.Register(gs)
	gs.Serve(l)
}
