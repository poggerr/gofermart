package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `db:"id"`
	Username string    `json:"login" db:"username"`
	Password string    `json:"password" db:"password"`
	Balance  float32   `json:"balance" db:"balance"`
	Withdraw float32   `json:"withdraw" db:"withdraw"`
}
