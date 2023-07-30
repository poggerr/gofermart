package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/poggerr/gophermart/internal/app"
	"github.com/poggerr/gophermart/internal/gzip"
	"github.com/poggerr/gophermart/internal/logger"
)

func Router(app *app.App) chi.Router {
	r := chi.NewRouter()
	r.Use(logger.WithLogging, gzip.GzipMiddleware)
	r.Post("/api/user/register", app.RegisterUser)
	r.Post("/api/user/login", app.UserLogin)
	r.Post("/api/user/orders", app.UploadOrder)
	r.Get("/api/user/orders", app.ListUserOrders)
	r.Get("/api/user/balance", app.UserBalance)
	r.Post("/api/user/balance/withdraw", app.Withdraw)
	r.Get("/api/user/withdrawals", app.Withdrawals)
	return r
}
