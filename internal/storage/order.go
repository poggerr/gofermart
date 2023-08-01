package storage

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/poggerr/gophermart/internal/accrual_service"
	"github.com/poggerr/gophermart/internal/logger"
	"github.com/poggerr/gophermart/internal/models"
	"strconv"
	"time"
)

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

type SaveOrd struct {
	OrderNum   string
	User       *uuid.UUID
	AccrualURL string
}

func (strg *Storage) SaveOrder(orderNumber int, user *uuid.UUID) error {
	t := time.Now()
	t.Format(time.RFC3339)
	id := uuid.New()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := strg.db.ExecContext(
		ctx,
		"INSERT INTO orders (id, order_number, order_user, uploaded_at, status, accrual_service) VALUES ($1, $2, $3, $4, $5, $6)",
		id, orderNumber, &user, t, "NEW", 0)
	if err != nil {
		logger.Initialize().Info(err)
		return err
	}
	return nil
}

func (strg *Storage) UpdateOrder(order SaveOrd) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	accrual, err := accrual_service.AccrualFun(order.OrderNum, order.AccrualURL)
	if err != nil {
		logger.Initialize().Info(err)
	}

	balance, err := strg.TakeUserBalance(order.User)
	if err != nil {
		logger.Initialize().Info(err)
	}
	if balance != nil {
		balance.Current += accrual.Accrual
		err = strg.UpdateUserBalance(order.User, balance.Current)
		if err != nil {
			logger.Initialize().Info(err)
		}

		fmt.Println(balance)
	}

	orderNumber, err := strconv.Atoi(order.OrderNum)
	if err != nil {
		logger.Initialize().Info(err)
	}

	_, err = strg.db.ExecContext(
		ctx,
		"UPDATE orders SET accrual_service=$1, status=$2 WHERE order_number=$3", accrual.Accrual, accrual.Status, orderNumber)
	if err != nil {
		logger.Initialize().Info(err)
	}
}
