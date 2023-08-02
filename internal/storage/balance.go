package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/poggerr/gophermart/internal/logger"
	"github.com/poggerr/gophermart/internal/models"
	"time"
)

func (strg *Storage) TakeUserBalance(userID *uuid.UUID) (*models.UserBalance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var userBalance models.UserBalance

	ans := strg.DB.QueryRowContext(ctx, "SELECT balance, withdrawn FROM main_user WHERE id=$1", userID)
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

	var balance models.UserBalance

	ans := strg.DB.QueryRowContext(ctx, "SELECT balance, withdrawn FROM main_user WHERE id=$1", userID)
	errScan := ans.Scan(&balance.Current, &balance.Withdrawn)
	if errScan != nil {
		logger.Initialize().Info(errScan)
		return errScan
	}
	balance.Current -= sum
	balance.Withdrawn += sum

	_, err := strg.DB.ExecContext(ctx, "UPDATE main_user SET balance=$1, withdrawn=$2 WHERE id=$3", balance.Current, balance.Withdrawn, userID)
	if err != nil {
		logger.Initialize().Info(err)
		return err
	}
	return nil
}

func (strg *Storage) UpdateUserBalance(userID *uuid.UUID, balance float32) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := strg.DB.ExecContext(
		ctx,
		"UPDATE main_user SET balance=$1 WHERE id=$2", balance, userID)
	if err != nil {
		return err
	}
	return nil
}
