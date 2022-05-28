package order

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"l0/pkg/domain"
	"log"
	"net/http"
)

type OUsecase interface {
	GetOrderByID(id int) (*domain.Order, error)
	AddOrder(order *domain.Order) error
}

type API struct {
	ou OUsecase
}

func NewAPI(ou OUsecase) *API {
	return &API{ou}
}

func (a *API) Order(w http.ResponseWriter, r *http.Request) {
	order, err := a.ou.GetOrderByID(1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(order.TrackNumber))
}

func (a *API) SubscribeToOrders(m *nats.Msg) {
	ord, err := orderFromJSON(m.Data)
	if err != nil {
		log.Println("error:", err)
		return
	}

	if err = a.ou.AddOrder(ord); err != nil {
		log.Fatal(err)
	}
}

func orderFromJSON(data []byte) (*domain.Order, error) {
	order := &domain.Order{}

	if err := json.Unmarshal(data, order); err != nil {
		return nil, fmt.Errorf("error unmarshalling order: %w", err)
	}

	return order, nil
}
