package accrual_service

import (
	"encoding/json"
	"github.com/poggerr/gophermart/internal/logger"
	"net/http"
	"time"
)

type Accrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual_service"`
}

func TakeAccrual(orderNumber string, url string) error {

	response, err := http.Get(url + "/api/orders/" + orderNumber)
	if err != nil {
		logger.Initialize().Info(err)
	}

	var ans Accrual

	dec := json.NewDecoder(response.Body)

	err = dec.Decode(&ans)
	if err != nil {
		logger.Initialize().Info(err)
	}

	return nil
}

func AccrualFun(orderNumber string, url string) error {
	//operation := TakeAccrual(orderNumber, url)

	//operation := func() error {
	//
	//}
	//err := backoff.Retry(operation, backoff.NewExponentialBackOff())
	//if err != nil {
	//	// Handle error.
	//	return
	//}

	for {
		err := TakeAccrual(orderNumber, url)
		if err != nil {
			logger.Initialize().Info(err)
			time.Sleep(1 * time.Millisecond)
		}
		break

	}
	return nil
}
