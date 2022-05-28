package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/nats-io/nats.go"
	"l0/configs"
	"l0/pkg/order"
	"log"
	"net/http"
	"os"
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

	_, err = nc.Subscribe(config.Topic, api.SubscribeToOrders)
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
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
}
