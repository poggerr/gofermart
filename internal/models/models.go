package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id       uuid.UUID `db:"id"`
	Username string    `json:"login" db:"username"`
	Password string    `json:"password" db:"password"`
}

type UserBalance struct {
	Current   float32 `json:"current" db:"balance"`
	Withdrawn int     `json:"withdrawn" db:"withdrawn"`
}

type UserOrder struct {
	Number     int       `db:"order_number" json:"number"`
	Status     string    `db:"status" json:"status"`
	Accrual    int       `db:"accrual" json:"accrual"`
	UploadedAt time.Time `db:"uploaded_at" json:"uploaded_at"`
}

type Orders []UserOrder

type Withdraw struct {
	OrderNumber string    `db:"order_number" json:"order"`
	Sum         float32   `db:"sum" json:"sum"`
	ProcessedAt time.Time `db:"processed_at" json:"processed_at"`
}

type Withdrawals []Withdraw
