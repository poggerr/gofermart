package app

import (
	"encoding/json"
	"github.com/poggerr/gophermart/internal/authorization"
	"github.com/poggerr/gophermart/internal/config"
	"github.com/poggerr/gophermart/internal/logger"
	"github.com/poggerr/gophermart/internal/models"
	"github.com/poggerr/gophermart/internal/service"
	"github.com/poggerr/gophermart/internal/storage"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"time"
)

type App struct {
	cfg           *config.Config
	strg          *storage.Storage
	sugaredLogger *zap.SugaredLogger
}

func NewApp(cfg *config.Config, strg *storage.Storage, log *zap.SugaredLogger) *App {
	return &App{
		cfg:           cfg,
		strg:          strg,
		sugaredLogger: log,
	}
}

func (a *App) RegisterUser(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		a.sugaredLogger.Info(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var user models.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		a.sugaredLogger.Info(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	isVerify := a.strg.VerifyUser(user.Username)
	if !isVerify {
		res.WriteHeader(http.StatusConflict)
		return
	}

	userID, err := authorization.RegisterUser(a.strg, &user)
	if err != nil {
		a.sugaredLogger.Info(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	jwtString, err := authorization.BuildJWTString(&userID)
	if err != nil {
		logger.Initialize().Info(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	cook := &http.Cookie{
		Name:    "session_token",
		Value:   jwtString,
		Path:    "/",
		Domain:  "localhost",
		Expires: time.Now().Add(120 * time.Second),
	}

	http.SetCookie(res, cook)

	res.WriteHeader(http.StatusOK)

}

func (a *App) UserLogin(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		a.sugaredLogger.Info(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var user models.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		a.sugaredLogger.Info(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, isVerify := a.strg.TakeUserID(user.Username)
	if isVerify {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	jwtString, err := authorization.BuildJWTString(userID)
	if err != nil {
		logger.Initialize().Info(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	cook := &http.Cookie{
		Name:    "session_token",
		Value:   jwtString,
		Path:    "/",
		Domain:  "localhost",
		Expires: time.Now().Add(120 * time.Second),
	}

	http.SetCookie(res, cook)

	res.WriteHeader(http.StatusOK)

}

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

	isValid := service.Valid(order)
	if !isValid {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	user, isUsed := a.strg.TakeOrderByUser(order)
	if isUsed {
		switch user {
		case userID:
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

	res.WriteHeader(http.StatusAccepted)

}

func (a *App) ListUserOrders(res http.ResponseWriter, req *http.Request) {
	//c, err := req.Cookie("session_token")
	//if err != nil {
	//	a.sugaredLogger.Info(err)
	//	res.WriteHeader(http.StatusUnauthorized)
	//	return
	//}

}

func (a *App) UserBalance(res http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("session_token")
	if err != nil {
		a.sugaredLogger.Info(err)
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID := authorization.GetUserID(c.Value)

	balance, err := a.strg.TakeUserBalance(userID)
	if err != nil {
		a.sugaredLogger.Info(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	marshal, err := json.Marshal(balance)
	if err != nil {
		a.sugaredLogger.Info(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json ")
	res.WriteHeader(http.StatusOK)
	res.Write(marshal)
}

type Withdraw struct {
	order int
	sum   float32
}

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

	var withdraw Withdraw

	err = json.Unmarshal(body, &withdraw)
	if err != nil {
		a.sugaredLogger.Info(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	_, isStore := a.strg.TakeOrderByUser(withdraw.order)
	if !isStore {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	balance, err := a.strg.TakeUserBalance(userID)
	if err != nil {
		a.sugaredLogger.Info(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if balance.Current < withdraw.sum {
		res.WriteHeader(http.StatusPaymentRequired)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func (a *App) Withdrawals(res http.ResponseWriter, req *http.Request) {

}
