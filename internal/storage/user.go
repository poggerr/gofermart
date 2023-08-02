package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/poggerr/gophermart/internal/models"
	"time"
)

func (strg *Storage) CreateUser(username, pass string, id *uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := strg.Db.ExecContext(
		ctx,
		"INSERT INTO main_user (id, username, password, withdrawn, balance) VALUES ($1, $2, $3, $4, $5)",
		id, username, pass, 0, 0)
	if err != nil {
		return err
	}
	return nil
}

func (strg *Storage) GetUser(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User

	ans := strg.Db.QueryRowContext(ctx, "SELECT * FROM main_user WHERE username=$1", username)
	errScan := ans.Scan(&user.ID, &user.Username, &user.Password, &user.Balance, &user.Withdraw)
	if errScan != nil {
		return nil, errScan
	}
	return &user, nil
}
