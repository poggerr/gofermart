package app

import (
	"github.com/poggerr/gophermart/internal/async"
	"github.com/poggerr/gophermart/internal/config"
	"github.com/poggerr/gophermart/internal/storage"
	"go.uber.org/zap"
)

type App struct {
	cfg           *config.Config
	strg          *storage.Storage
	sugaredLogger *zap.SugaredLogger
	repo          *async.AccrualRepo
}

func NewApp(cfg *config.Config, strg *storage.Storage, log *zap.SugaredLogger, repo *async.AccrualRepo) *App {
	return &App{
		cfg:           cfg,
		strg:          strg,
		sugaredLogger: log,
		repo:          repo,
	}
}
