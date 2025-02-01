package models

type Wallet struct {
	ID 	     int64  `json:"id"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
	Owner    int64  `json:"owner"`
}