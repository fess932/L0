package main

import (
	"l0/configs"
	"l0/pkg/order"
	"log"
	"net/http"
)

func main() {
	config := configs.GetConfig()

	repo := order.NewPgRepo(nil, configs.GetConfig())
	ucase := order.NewUsecase(repo)
	api := order.NewAPI(ucase)

	mux := http.NewServeMux()
	mux.HandleFunc("/", api.Order)

	log.Println("server listen on", config.Host)
	log.Println(http.ListenAndServe(config.Host, mux))
}

func setupPg() {
	//TODO implement me
}

func connNatsStreaming() {
	//TODO implement me
}
