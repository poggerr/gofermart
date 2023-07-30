package models

type UserBalance struct {
	Current   float32 `json:"current" db:"balance"`
	Withdrawn int     `json:"withdrawn" db:"withdrawn"`
}
