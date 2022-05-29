package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coocood/freecache"
	"github.com/jackc/pgx/v4"
	"l0/configs"
	"l0/pkg/domain"
	"log"
	"strconv"
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

func (p *PgRepo) GetOrderByID(ctx context.Context, id int) (order *domain.Order, err error) {
	val, err := p.cache.Get([]byte(strconv.Itoa(id)))

	switch {
	case errors.Is(err, freecache.ErrNotFound):
		log.Println("order not found in cache")
	case err != nil:
		log.Println("error getting order from cache: %w", err)
	default:
		if err = json.Unmarshal(val, &order); err != nil {
			log.Println("error unmarshalling order from cache: %w", err)

			break
		}

		return order, nil
	}

	if err = p.db.QueryRow(ctx, "SELECT * FROM orders WHERE id = $1", id).
		Scan(&id, order); err != nil {
		return nil, fmt.Errorf("error getting order: %w", err)
	}

	return order, nil
}

func (p *PgRepo) AddOrder(ctx context.Context, order *domain.Order) error {
	err := p.db.QueryRow(ctx,
		`
INSERT INTO orders (order_uid, track_number, entry, date_created) 
VALUES ($1, $2, $3, $4) 
RETURNING id`,
		order.OrderUID, order.TrackNumber, order.Entry, order.DateCreated,
	).Scan(&order.ID)

	if err != nil {
		return fmt.Errorf("error adding order: %w", err)
	}

	o, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("error marshalling order: %w", err)
	}

	err = p.cache.Set([]byte(strconv.Itoa(order.ID)), o, 0)
	if err != nil {
		return fmt.Errorf("error adding order to cache: %w", err)
	}

	return nil
}
