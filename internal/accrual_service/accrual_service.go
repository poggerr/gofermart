package accrual_service

import (
	"encoding/json"
	"github.com/poggerr/gophermart/internal/logger"
	"net/http"
)

type Accrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual_service"`
}

func TakeAccrual(orderNumber string, url string) (float32, error) {

	response, err := http.Get(url + "/api/orders/" + orderNumber)
	if err != nil {
		logger.Initialize().Info(err)
		return 0, err
	}

	var ans Accrual

	dec := json.NewDecoder(response.Body)

	err = dec.Decode(&ans)
	if err != nil {
		logger.Initialize().Info(err)
		return 0, err
	}

	return ans.Accrual, nil
}
