package app

import (
	"encoding/json"
	"fmt"
	"github.com/poggerr/gophermart/internal/authorization"
	"net/http"
)

func (a *App) ListUserOrders(res http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("session_token")
	if err != nil {
		a.sugaredLogger.Info(err)
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID := authorization.GetUserID(c.Value)

	orders, err := a.strg.TakeUserOrders(userID)
	if err != nil {
		a.sugaredLogger.Info(err)
		res.WriteHeader(http.StatusNoContent)
		return
	}

	fmt.Println(orders)

	marshal, err := json.Marshal(orders)
	if err != nil {
		a.sugaredLogger.Info(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json ")
	res.WriteHeader(http.StatusOK)
	res.Write(marshal)

}
