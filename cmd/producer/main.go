package main

import (
	"encoding/json"
	"l0/pkg/domain"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	stan "github.com/nats-io/stan.go"
	"l0/configs"
)

func main() {
	sc, err := stan.Connect(
		configs.GetConfig().NatsClusterID,
		"order_producer",
		stan.NatsURL(configs.GetConfig().NatsHost),
	)
	if err != nil {
		log.Println(err)

		return
	}
	defer sc.Close()

	t := time.NewTicker(time.Nanosecond * 1)
	for range t.C {
		log.Println("send msg to topic", configs.GetConfig().Topic)

		if err = sc.Publish(configs.GetConfig().Topic, genOrder()); err != nil {
			log.Fatal(err)
		}
	}
}

func genOrder() []byte {
	order := domain.Order{
		SmID:              gofakeit.Number(1, 100),
		OrderUID:          gofakeit.UUID(),
		TrackNumber:       gofakeit.UUID(),
		Entry:             "WBIL",
		Locale:            gofakeit.Language(),
		InternalSignature: "",
		CustomerID:        gofakeit.UUID(),
		DeliveryService:   "meest",
		Shardkey:          "9",
		OofShard:          "1",
		DateCreated:       gofakeit.Date(),

		Items: []domain.Item{
			{
				ChrtID:      gofakeit.IntRange(1, 10000),
				TrackNumber: gofakeit.UUID(),
				Price:       gofakeit.IntRange(1, 10000),
				Rid:         gofakeit.UUID(),
				Name:        gofakeit.BeerName(),
				Sale:        0,
				Size:        "",
				TotalPrice:  0,
				NmID:        0,
				Brand:       "",
				Status:      0,
			},
		},
		Delivery: domain.Delivery{
			Name:    gofakeit.StreetName(),
			Phone:   gofakeit.Phone(),
			Zip:     gofakeit.Zip(),
			City:    gofakeit.City(),
			Address: gofakeit.BitcoinAddress(),
			Region:  "",
			Email:   gofakeit.Email(),
		},
		Payment: domain.Payment{
			Transaction:  gofakeit.UUID(),
			RequestID:    "",
			Currency:     "",
			Provider:     "",
			Amount:       0,
			PaymentDt:    0,
			Bank:         "",
			DeliveryCost: 0,
			GoodsTotal:   0,
			CustomFee:    0,
		},
	}

	b, err := json.Marshal(order)
	if err != nil {
		log.Fatal(err)
	}

	return b
}
