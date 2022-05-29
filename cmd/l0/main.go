package main

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/nats-io/nats.go"
	"l0/configs"
	"l0/pkg/order"
	"log"
	"net/http"
)

func main() {
	config := configs.GetConfig()
	db := pgConn()

	repo := order.NewPgRepo(db)
	if err := repo.RefreshCacheFromDB(context.Background(), config); err != nil {
		log.Println("error refreshing cache from db", err)
	}

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
	mux.HandleFunc("/", api.InputOrderIDHandler)
	log.Println("server listen on", config.Host)
	log.Println(http.ListenAndServe(config.Host, mux))
}

func pgConn() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), configs.GetConfig().DB)
	if err != nil {
		log.Fatal(err)
	}

	if err = conn.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}

	return conn
}
