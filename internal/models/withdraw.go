package models

import "time"

type Withdraw struct {
	OrderNumber string    `db:"order_number" json:"order"`
	Sum         float32   `db:"sum" json:"sum"`
	ProcessedAt time.Time `db:"processed_at" json:"processed_at"`
}

type Withdrawals []Withdraw
