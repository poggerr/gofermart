package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/poggerr/gophermart/internal/app"
	"github.com/poggerr/gophermart/internal/authorization"
	"github.com/poggerr/gophermart/internal/gzip"
	"github.com/poggerr/gophermart/internal/logger"
)

func Router(app *app.App) chi.Router {
	r := chi.NewRouter()
	r.Use(logger.WithLogging, gzip.GzipMiddleware)
	r.Post("/api/user/register", app.RegisterUser)
	r.Post("/api/user/login", app.UserLogin)
	r.With(authorization.AuthMiddleware).Post("/api/user/orders", app.UploadOrder)
	r.With(authorization.AuthMiddleware).Get("/api/user/orders", app.ListUserOrders)
	r.With(authorization.AuthMiddleware).Get("/api/user/balance", app.UserBalance)
	r.With(authorization.AuthMiddleware).Post("/api/user/balance/withdraw", app.Withdraw)
	r.With(authorization.AuthMiddleware).Get("/api/user/withdrawals", app.Withdrawals)
	return r
}
