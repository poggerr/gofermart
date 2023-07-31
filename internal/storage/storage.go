package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/poggerr/gophermart/internal/accrual_service"
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
CREATE TABLE IF NOT EXISTS main_user (
    id UUID UNIQUE,
    username text,
    password text,
    balance float,
    withdrawn int
);

CREATE TABLE IF NOT EXISTS orders (
    id UUID UNIQUE,
    order_number bigint UNIQUE,
    order_user UUID,
    uploaded_at date,
    accrual_service float,
    status text
);

CREATE TABLE IF NOT EXISTS withdrawals (
    id UUID UNIQUE,
    order_number bigint UNIQUE,
    order_user UUID,
    sum float,
    processed_at date
)
`

func (strg *Storage) RestoreDB() {
	strg.db.MustExec(schema)
}

func (strg *Storage) CreateUser(username, pass string, id *uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := strg.db.ExecContext(
		ctx,
		"INSERT INTO main_user (id, username, password, withdrawn, balance) VALUES ($1, $2, $3, $4, $5)",
		id, username, pass, 0, 0)
	if err != nil {
		logger.Initialize().Info(err)
	}
	return nil
}

func (strg *Storage) VerifyUser(username string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ans := strg.db.QueryRowContext(ctx, "SELECT username FROM main_user WHERE username=$1", username)
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

	ans := strg.db.QueryRowContext(ctx, "SELECT id FROM main_user WHERE username=$1", username)
	errScan := ans.Scan(&id)
	if errScan != nil {
		logger.Initialize().Info(errScan)
		return nil, true
	}
	return &id, false
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

	ans := strg.db.QueryRowContext(ctx, "SELECT balance, withdrawn FROM main_user WHERE id=$1", userID)
	errScan := ans.Scan(&userBalance.Current, &userBalance.Withdrawn)
	if errScan != nil {
		logger.Initialize().Info(errScan)
		return nil, errScan
	}
	return &userBalance, nil
}

func (strg *Storage) Debit(userID *uuid.UUID, sum float32) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var oldSum float32

	ans := strg.db.QueryRowContext(ctx, "SELECT balance FROM main_user WHERE id=$1", userID)
	errScan := ans.Scan(&oldSum)
	if errScan != nil {
		logger.Initialize().Info(errScan)
		return errScan
	}
	oldSum -= sum
	_, err := strg.db.ExecContext(ctx, "UPDATE main_user SET balance=$1 WHERE id=$2", oldSum, userID)
	if err != nil {
		logger.Initialize().Info(err)
		return err
	}
	return nil
}

func (strg *Storage) CreateWithdraw(userID *uuid.UUID, withdraw *models.Withdraw) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	id := uuid.New()
	t := time.Now()
	t.Format(time.RFC3339)

	_, err := strg.db.ExecContext(
		ctx,
		"INSERT INTO withdrawals (id, order_number, order_user, sum, processed_at) VALUES ($1, $2, $3, $4, $5)",
		id, withdraw.OrderNumber, userID, withdraw.Sum, t)
	if err != nil {
		logger.Initialize().Info(err)
		return err
	}
	return nil
}

func (strg *Storage) TakeUserOrders(userID *uuid.UUID) (*models.Orders, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := strg.db.QueryContext(ctx, "SELECT * FROM orders WHERE order_user=$1", userID)
	if err != nil {
		logger.Initialize().Info(err)
		return nil, err
	}

	orders := make(models.Orders, 0)
	for rows.Next() {
		var order models.UserOrder
		var id uuid.UUID
		var orderUser uuid.UUID
		if err = rows.Scan(&id, &order.Number, &orderUser, &order.UploadedAt, &order.Accrual, &order.Status); err != nil {
			logger.Initialize().Info(err)
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		logger.Initialize().Info(err)
		return nil, err
	}
	return &orders, nil
}

func (strg *Storage) TakeUserWithdrawals(userID *uuid.UUID) (*models.Withdrawals, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := strg.db.QueryContext(ctx, "SELECT * FROM withdrawals WHERE order_user=$1", userID)
	if err != nil {
		logger.Initialize().Info(err)
		return nil, err
	}

	withdrawals := make(models.Withdrawals, 0)
	for rows.Next() {
		var withdraw models.Withdraw
		var id uuid.UUID
		var orderUser uuid.UUID
		if err = rows.Scan(&id, &withdraw.OrderNumber, &orderUser, &withdraw.Sum, &withdraw.ProcessedAt); err != nil {
			logger.Initialize().Info(err)
			return nil, err
		}
		withdrawals = append(withdrawals, withdraw)
	}

	if err = rows.Err(); err != nil {
		logger.Initialize().Info(err)
		return nil, err
	}
	return &withdrawals, nil
}

func (strg *Storage) TakeUserPass(user *models.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var dbPass string

	ans := strg.db.QueryRowContext(ctx, "SELECT password FROM main_user WHERE username=$1", user.Username)
	errScan := ans.Scan(&dbPass)
	if errScan != nil {
		logger.Initialize().Info(errScan)
		return "", errScan
	}
	return dbPass, nil
}

func (strg *Storage) UpdateUserBalance(userID *uuid.UUID, balance float32) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := strg.db.ExecContext(
		ctx,
		"UPDATE main_user SET balance=$1 WHERE id=$2", balance, userID)
	if err != nil {
		return err
	}
	return nil
}

type SaveOrd struct {
	OrderNum   int
	User       *uuid.UUID
	AccrualURL string
}

func (strg *Storage) SaveOrder(order SaveOrd) {
	t := time.Now()
	t.Format(time.RFC3339)
	id := uuid.New()

	accrual, err := accrual_service.AccrualFun(string(rune(order.OrderNum)), order.AccrualURL)
	if err != nil {
		logger.Initialize().Info(err)
	}

	balance, err := strg.TakeUserBalance(order.User)
	if err != nil {
		logger.Initialize().Info(err)
	}

	balance.Current += accrual

	err = strg.UpdateUserBalance(order.User, balance.Current)
	if err != nil {
		logger.Initialize().Info(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = strg.db.ExecContext(
		ctx,
		"INSERT INTO orders (id, order_number, order_user, uploaded_at, status, accrual_service) VALUES ($1, $2, $3, $4, $5, $6)",
		id, order.OrderNum, &order.User, t, "NEW", accrual)
	if err != nil {
		logger.Initialize().Info(err)
	}
}
