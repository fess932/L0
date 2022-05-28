package main

import (
	"github.com/nats-io/nats.go"
	"l0/configs"
	"log"
	"time"
)

func main() {
	nc, err := nats.Connect(configs.GetConfig().NatsHost)
	if err != nil {
		log.Fatal(err)
	}

	t := time.NewTicker(time.Second * 1)
	for range t.C {
		log.Println("send msg to topic", configs.GetConfig().Topic)
		if err = nc.Publish(configs.GetConfig().Topic, []byte("Hello World")); err != nil {
			log.Fatal(err)
		}
	}
}
