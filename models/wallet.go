package models

type Wallet struct {
	Id       int64   `json:"id" field:"id"`
	OwnerId  int64   `json:"ownerId" field:"ownerid"`
	Name     string  `json:"name" validate:"omitempty,max=50" field:"name"`
	Currency string  `json:"currency" validate:"omitempty,iso4217" field:"currency"`
	Balance  float64 `json:"balance" field:"balance"`
}
