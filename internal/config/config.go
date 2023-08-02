package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v8"
	"github.com/cenkalti/backoff/v4"
	"net/http"
	"time"
)

type Config struct {
	ServAddr string `env:"RUN_ADDRESS"`
	DB       string `env:"DATABASE_URI"`
	Accrual  string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	Client   *http.Client
	Backoff  *backoff.ExponentialBackOff
}

func NewConf() *Config {
	var cfg Config

	flag.StringVar(&cfg.ServAddr, "a", ":8080", "write down server")
	flag.StringVar(
		&cfg.DB,
		"d",
		"host=localhost user=gophermart password=userpassword dbname=gophermart sslmode=disable",
		"write down db")
	flag.StringVar(&cfg.Accrual, "r", "http://localhost:8080", "write down accrual_service server")
	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	cfg.Client = &http.Client{}
	cfg.Backoff = backoff.NewExponentialBackOff()
	cfg.Backoff.MaxElapsedTime = 10 * time.Second

	return &cfg
}
