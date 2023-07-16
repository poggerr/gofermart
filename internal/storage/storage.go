package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/poggerr/gophermart/internal/logger"
	"github.com/poggerr/gophermart/internal/models"
	"time"
)

type Storage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		db: db,
	}
}

var schema = `
CREATE TABLE IF NOT EXISTS MainUser (
    id UUID UNIQUE,
    username text,
    password text,
    balance float,
    withdrawn int
);

CREATE TABLE IF NOT EXISTS orders (
    id UUID UNIQUE,
    order_number int UNIQUE,
    order_user UUID,
    uploaded_at date
)
`

func (strg *Storage) RestoreDB() {
	strg.db.MustExec(schema)
}

func (strg *Storage) CreateUser(username, pass string, id *uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := strg.db.ExecContext(ctx, "INSERT INTO mainuser (id, username, password) VALUES ($1, $2, $3)", id, username, pass)
	if err != nil {
		logger.Initialize().Info(err)
	}
	return nil
}

func (strg *Storage) VerifyUser(username string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ans := strg.db.QueryRowContext(ctx, "SELECT username FROM mainuser WHERE username=$1", username)
	errScan := ans.Scan(&username)
	if errScan != nil {
		logger.Initialize().Info(errScan)
		return true
	}
	return false
}

func (strg *Storage) TakeUserID(username string) (*uuid.UUID, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id uuid.UUID

	ans := strg.db.QueryRowContext(ctx, "SELECT id FROM mainuser WHERE username=$1", username)
	errScan := ans.Scan(&id)
	if errScan != nil {
		logger.Initialize().Info(errScan)
		return nil, true
	}
	return &id, false
}

func (strg *Storage) SaveOrder(orderNumber int, user *uuid.UUID) error {
	t := time.Now()
	t.Format(time.RFC3339)

	id := uuid.New()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := strg.db.ExecContext(ctx, "INSERT INTO orders (id, order_number, order_user, uploaded_at) VALUES ($1, $2, $3, $4)", id, orderNumber, &user, t)
	if err != nil {
		logger.Initialize().Info(err)
		return err
	}
	return nil
}

func (strg *Storage) TakeOrderByUser(orderNumber int) (*uuid.UUID, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user uuid.UUID

	ans := strg.db.QueryRowContext(ctx, "SELECT order_user FROM orders WHERE order_number=$1", orderNumber)
	errScan := ans.Scan(&user)
	if errScan != nil {
		logger.Initialize().Info(errScan)
		return nil, false
	}
	return &user, true
}

func (strg *Storage) TakeUserBalance(userID *uuid.UUID) (*models.UserBalance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var userBalance models.UserBalance

	ans := strg.db.QueryRowContext(ctx, "SELECT balance, withdrawn FROM mainuser WHERE id=$1", userID)
	errScan := ans.Scan(&userBalance)
	if errScan != nil {
		logger.Initialize().Info(errScan)
		return nil, errScan
	}
	return &userBalance, nil
}
