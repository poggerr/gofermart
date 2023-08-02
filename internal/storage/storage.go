package storage

import (
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	Db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		Db: db,
	}
}
