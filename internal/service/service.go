package service

import (
	"encoding/json"
	"fmt"
	"github.com/poggerr/gophermart/internal/logger"
	"io"
	"net/http"
)

func Valid(number int) bool {
	return (number%10+checksum(number/10))%10 == 0
}

func checksum(number int) int {
	var luhn int
	for i := 0; number > 0; i++ {
		cur := number % 10
		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}
		luhn += cur
		number = number / 10
	}
	return luhn % 10
}

type Accrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual"`
}

func TakeAccrual(orderNumber string, url string) (float32, error) {
	response, err := http.Get(url + "/api/orders/" + orderNumber)
	if err != nil {
		logger.Initialize().Info(err)
		return 0, err
	}

	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		logger.Initialize().Info(err)
		return 0, err
	}

	var ans Accrual

	fmt.Println(string(body))

	err = json.Unmarshal(body, &ans)
	if err != nil {
		logger.Initialize().Info(err)
		return 0, err
	}

	fmt.Println(ans)

	return ans.Accrual, nil
}
