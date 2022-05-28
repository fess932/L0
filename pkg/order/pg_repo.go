package order

import (
	"context"
	"fmt"
	"github.com/coocood/freecache"
	"github.com/jackc/pgx/v4"
	"l0/configs"
	"l0/pkg/domain"
	"log"
)

type PgRepo struct {
	db    *pgx.Conn
	cache *freecache.Cache
}

func NewPgRepo(db *pgx.Conn, conf configs.Config) *PgRepo {
	repo := &PgRepo{
		db:    db,
		cache: freecache.NewCache(conf.CacheSize),
	}

	return repo
}

func (p *PgRepo) GetOrderByID(id int) (order *domain.Order, err error) {
	if err = p.db.QueryRow(context.Background(), "SELECT * FROM orders WHERE id = $1", id).
		Scan(&id, order); err != nil {
		return nil, fmt.Errorf("error getting order: %w", err)
	}

	return order, nil
}

func (p *PgRepo) AddOrder(order *domain.Order) error {
	t, err := p.db.Exec(context.Background(),
		"INSERT INTO orders (order_uid, track_number, entry, date_created) VALUES ($1, $2, $3, $4)",
		order.OrderUID, order.TrackNumber, order.Entry, order.DateCreated,
	)

	if err != nil {
		return fmt.Errorf("error adding order: %w", err)
	}

	log.Println("tag", t.String(), t.RowsAffected())

	return nil
}
