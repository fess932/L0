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

func NewPgRepo(db *pgx.Conn) *PgRepo {
	repo := &PgRepo{
		db: db,
	}

	return repo
}

func scanOrder(o *domain.Order, row pgx.Row) (err error) {
	if err = row.Scan(&o.OrderUID, &o.TrackNumber, &o.Entry, &o.DateCreated); err != nil {
		return fmt.Errorf("error scanning order id: %w", err)
	}

	return nil
}

func (p *PgRepo) RefreshCacheFromDB(ctx context.Context, conf configs.Config) error {
	const avgOrderSize = 2 * 1024 // average size of order is 2kb(model.json)
	maxCacheEntries := conf.CacheSize / avgOrderSize
	p.cache = freecache.NewCache(conf.CacheSize)

	rows, err := p.db.Query(ctx, `
SELECT order_uid, track_number, entry, date_created 
FROM orders 
ORDER BY id DESC
LIMIT $1`, maxCacheEntries)
	if err != nil {
		return fmt.Errorf("error getting orders from db: %w", err)
	}
	defer rows.Close()

	var (
		order domain.Order
		data  []byte
	)

	for rows.Next() {
		order = domain.Order{}
		if err = scanOrder(&order, rows); err != nil {
			return fmt.Errorf("error fetching order: %w", err)
		}

		data, err = json.Marshal(order)
		if err != nil {
			return fmt.Errorf("error marshalling order: %w", err)
		}

		err = p.cache.Set([]byte(strconv.Itoa(order.ID)), data, 0)
		if err != nil {
			return fmt.Errorf("error adding order to cache: %w", err)
		}
	}

	log.Println("cache refreshed")

	return nil
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

	order = &domain.Order{ID: id}

	err = scanOrder(order, p.db.QueryRow(ctx, `
SELECT order_uid, track_number, entry, date_created
FROM orders 
WHERE id=$1
`, id))

	if err != nil {
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
