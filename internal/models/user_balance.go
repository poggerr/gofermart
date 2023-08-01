package models

type UserBalance struct {
	Current   float32 `json:"current" db:"balance"`
	Withdrawn float32 `json:"withdrawn" db:"withdrawn"`
}
