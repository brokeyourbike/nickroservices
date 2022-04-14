package main

import (
	"net/http"

	"github.com/brokeyourbike/nickroservices/handlers"
	"github.com/brokeyourbike/nickroservices/protos"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
)

func main() {
	log := hclog.Default()

	conn, err := grpc.Dial("127.0.0.1:9092", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	cc := protos.NewCurrencyClient(conn)

	ph := handlers.NewProducts(log, cc)

	r := mux.NewRouter()
	r.HandleFunc("/", ph.GetProduct)

	s := http.Server{
		Addr:    "127.0.0.1:9090",
		Handler: r,
	}

	err = s.ListenAndServe()
	if err != nil {
		log.Error("Cannot run server", err)
	}
}
