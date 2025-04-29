package models

type Wallet struct {
	ID       int64   `json:"id" field:"id"`
	OwnerID  int64   `json:"owner_id" field:"ownerid"`
	Name     string  `json:"name" validate:"omitempty,max=50" field:"name"`
	Currency string  `json:"currency" validate:"omitempty,max=50" field:"currency"`
	Balance  float64 `json:"balance" field:"balance"`
}
