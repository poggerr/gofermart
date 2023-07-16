package models

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `db:"id"`
	Username string    `json:"login" db:"username"`
	Password string    `json:"password" db:"password"`
}

type UserBalance struct {
	Current   float32 `json:"current" db:"balance"`
	Withdrawn int     `json:"withdrawn" db:"withdrawn"`
}
