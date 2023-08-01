package app

import (
	"github.com/poggerr/gophermart/internal/authorization"
	"github.com/poggerr/gophermart/internal/orderValidation"
	"io"
	"net/http"
	"strconv"
)

func (a *App) UploadOrder(res http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("session_token")
	if err != nil {
		a.sugaredLogger.Info(err)
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID := authorization.GetUserID(c.Value)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		a.sugaredLogger.Info(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	order, err := strconv.Atoi(string(body))
	if err != nil {
		a.sugaredLogger.Info(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	isValid := orderValidation.OrderValidation(order)
	if !isValid {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	user, isUsed := a.strg.TakeOrderByUser(order)
	if isUsed {
		switch *user {
		case *userID:
			res.WriteHeader(http.StatusOK)
			return
		default:
			res.WriteHeader(http.StatusConflict)
			return
		}
	}

	err = a.strg.SaveOrder(order, userID)
	if err != nil {
		a.sugaredLogger.Info(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	a.repo.TakeAsync(string(body), userID, a.cfg.Accrual)

	res.WriteHeader(http.StatusAccepted)

}
