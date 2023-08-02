package app

import (
	"context"
	"github.com/google/uuid"
	"github.com/poggerr/gophermart/internal/logger"
	"github.com/poggerr/gophermart/internal/models"
	"time"
)

func (a *App) AccrualRestore() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := a.strg.DB.QueryContext(ctx, "SELECT * FROM orders WHERE status=$1 OR status=$2 OR status=$3", "NEW", "REGISTERED", "PROCESSING")
	if err != nil {
		logger.Initialize().Info(err)
		return
	}
	for rows.Next() {
		var order models.UserOrder
		var id uuid.UUID
		var orderUser uuid.UUID
		if err = rows.Scan(&id, &order.Number, &orderUser, &order.UploadedAt, &order.Accrual, &order.Status); err != nil {
			logger.Initialize().Info(err)
		}
		a.repo.SendToChan(order.Number, &orderUser, a.cfg.Accrual)
	}

	if err = rows.Err(); err != nil {
		logger.Initialize().Info(err)
	}
}
