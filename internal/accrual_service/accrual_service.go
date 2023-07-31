package accrual_service

import (
	"encoding/json"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/poggerr/gophermart/internal/logger"
	"net/http"
	"time"
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

func AccrualFun(orderNumber string, url string) (float32, error) {
	client := &http.Client{}
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 20 * time.Second

	var ans Accrual

	operation := func() error {
		resp, err := client.Get(url + "/api/orders/" + orderNumber)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		fmt.Println(resp.StatusCode)

		if resp.StatusCode == http.StatusOK {
			dec := json.NewDecoder(resp.Body)

			err = dec.Decode(&ans)
			if err != nil {
				logger.Initialize().Info(err)
			}
			return nil
		}

		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	err := backoff.Retry(operation, b)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return 0, err
	}

	return ans.Accrual, nil
}
