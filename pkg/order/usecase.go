package order

import (
	"context"
	"fmt"
	"l0/pkg/domain"
	"log"
)

type ORepo interface {
	GetOrderByID(ctx context.Context, id int) (*domain.Order, error)
	AddOrder(ctx context.Context, order *domain.Order) error
}

func NewUsecase(or ORepo) *Usecase {
	return &Usecase{or}
}

type Usecase struct {
	or ORepo
}

func (u *Usecase) AddOrder(order *domain.Order) error {
	if err := u.or.AddOrder(context.Background(), order); err != nil {
		return fmt.Errorf("error adding order: %w", err)
	}

	log.Printf("order added: %v\n", order.ID)

	return nil
}

func (u *Usecase) GetOrderByID(id int) (*domain.Order, error) {
	return u.or.GetOrderByID(context.Background(), id)
}
