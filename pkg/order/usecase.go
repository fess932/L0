package order

import (
	"l0/pkg/domain"
	"log"
)

type ORepo interface {
	GetOrderByID(id int) (*domain.Order, error)
}

type Usecase struct {
	or ORepo
}

func (u *Usecase) AddOrder(order *domain.Order) error {
	log.Println("inmplement me")

	return nil
}

func NewUsecase(or ORepo) *Usecase {
	return &Usecase{or}
}

func (u *Usecase) GetOrderByID(id int) (*domain.Order, error) {
	return u.or.GetOrderByID(id)
}
