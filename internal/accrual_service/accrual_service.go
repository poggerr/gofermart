package accrual_service

import (
	"encoding/json"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/poggerr/gophermart/internal/logger"
	"github.com/poggerr/gophermart/internal/models"
	"net/http"
	"time"
)

func AccrualFun(orderNumber string, url string) (*models.Accrual, error) {
	client := &http.Client{}
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 10 * time.Second

	var ans models.Accrual

	operation := func() error {
		resp, err := client.Get(url + "/api/orders/" + orderNumber)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		dec := json.NewDecoder(resp.Body)

		err = dec.Decode(&ans)
		if err != nil {
			logger.Initialize().Info(err)
		}

		if ans.Status == "PROCESSED" || ans.Status == "INVALID" {
			return nil
		}

		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	err := backoff.Retry(operation, b)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil, err
	}

	fmt.Println(ans)

	return &ans, nil
}
