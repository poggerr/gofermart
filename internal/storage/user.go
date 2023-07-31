package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/poggerr/gophermart/internal/logger"
	"github.com/poggerr/gophermart/internal/models"
	"time"
)

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
