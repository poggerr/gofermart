package app

import (
	"encoding/json"
	"github.com/poggerr/gophermart/internal/authorization"
	"github.com/poggerr/gophermart/internal/models"
	"io"
	"net/http"
	"strconv"
)

func (a *App) Withdraw(res http.ResponseWriter, req *http.Request) {
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

	var withdraw models.Withdraw

	err = json.Unmarshal(body, &withdraw)
	if err != nil {
		a.sugaredLogger.Info(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	order, err := strconv.Atoi(withdraw.OrderNumber)
	if err != nil {
		a.sugaredLogger.Info(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	_, isStore := a.strg.TakeOrderByUser(order)
	if isStore {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	balance, err := a.strg.TakeUserBalance(userID)
	if err != nil {
		a.sugaredLogger.Info(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if balance.Current < withdraw.Sum {
		res.WriteHeader(http.StatusPaymentRequired)
		return
	}

	err = a.strg.Debit(userID, withdraw.Sum)
	if err != nil {
		a.sugaredLogger.Info(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = a.strg.CreateWithdraw(userID, &withdraw)
	if err != nil {
		a.sugaredLogger.Info(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}
