package app

import (
	"github.com/poggerr/gophermart/internal/config"
	"github.com/poggerr/gophermart/internal/repo"
	"github.com/poggerr/gophermart/internal/storage"
	"go.uber.org/zap"
)

type App struct {
	cfg           *config.Config
	strg          *storage.Storage
	sugaredLogger *zap.SugaredLogger
	repo          *repo.AccrualRepo
}

func NewApp(cfg *config.Config, strg *storage.Storage, log *zap.SugaredLogger, repo *repo.AccrualRepo) *App {
	return &App{
		cfg:           cfg,
		strg:          strg,
		sugaredLogger: log,
		repo:          repo,
	}
}
