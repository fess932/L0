package main

import (
	"github.com/nats-io/nats.go"
	"l0/configs"
	"l0/pkg/order"
	"log"
	"net/http"
)

func main() {
	config := configs.GetConfig()
	repo := order.NewPgRepo(nil, config)
	ucase := order.NewUsecase(repo)
	api := order.NewAPI(ucase)

	// streaming
	nc, err := nats.Connect(config.NatsHost)
	if err != nil {
		log.Fatal(err)
	}

	_, err = nc.Subscribe(config.Topic, func(m *nats.Msg) {
		log.Println("INCOME:", string(m.Data))
	})
	if err != nil {
		log.Fatal(err)
	}

	// http
	mux := http.NewServeMux()
	mux.HandleFunc("/", api.Order)
	log.Println("server listen on", config.Host)
	log.Println(http.ListenAndServe(config.Host, mux))
}

func setupPg() {
	//TODO implement me
}
