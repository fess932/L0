package order

import (
	"database/sql"
	"github.com/coocood/freecache"
	"l0/configs"
	"l0/pkg/domain"
)

type PgRepo struct {
	db    *sql.DB
	cache *freecache.Cache
}

func (p PgRepo) GetOrderByID(id int) (*domain.Order, error) {
	//TODO implement me
	panic("implement me")
}

func NewPgRepo(db *sql.DB, conf configs.Config) *PgRepo {
	repo := &PgRepo{
		db:    db,
		cache: freecache.NewCache(conf.CacheSize),
	}

	return repo
}
