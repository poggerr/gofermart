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
