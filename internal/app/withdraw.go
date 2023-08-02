package app

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/poggerr/gophermart/internal/authorization"
	"github.com/poggerr/gophermart/internal/models"
	"io"
	"net/http"
	"strconv"
)

func (a *App) checkBalance(withdraw *models.Withdraw, userID *uuid.UUID) error {
	balance, err := a.strg.TakeUserBalance(userID)
	if err != nil {
		a.sugaredLogger.Info(err)
		return err
	}

	if balance.Current < withdraw.Sum {
		return err
	}
	return nil
}

func buildWithdraw(req *http.Request) (*models.Withdraw, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var withdraw models.Withdraw

	err = json.Unmarshal(body, &withdraw)
	if err != nil {
		return nil, err
	}

	return &withdraw, err
}

func (a *App) Withdraw(res http.ResponseWriter, req *http.Request) {
	userID := authorization.FromContext(req.Context())

	withdraw, err := buildWithdraw(req)
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

	err = a.checkBalance(withdraw, userID)
	if err != nil {
		res.WriteHeader(http.StatusPaymentRequired)
		return
	}

	err = a.strg.Debit(userID, withdraw.Sum)
	if err != nil {
		a.sugaredLogger.Info(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = a.strg.CreateWithdraw(userID, withdraw)
	if err != nil {
		a.sugaredLogger.Info(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}
