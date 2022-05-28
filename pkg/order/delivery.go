package order

import (
	"l0/pkg/domain"
	"net/http"
)

type OUsecase interface {
	GetOrderByID(id int) (*domain.Order, error)
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
