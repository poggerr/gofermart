package storage

import (
	"github.com/jmoiron/sqlx"
	"github.com/poggerr/gophermart/internal/config"
)

type Storage struct {
	DB  *sqlx.DB
	cfg *config.Config
}

func NewStorage(db *sqlx.DB, cfg *config.Config) *Storage {
	return &Storage{
		DB:  db,
		cfg: cfg,
	}
}
