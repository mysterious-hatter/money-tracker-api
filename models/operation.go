package models

type Operation struct {
	ID         int64   `json:"id" field:"id"`
	Name       string  `json:"name" validate:"omitempty,max=50" field:"name"`
	WalletID   int64   `json:"wallet_id" field:"walletid"`
	Sum        float64 `json:"sum" field:"sum"`
	Date       string  `json:"date" validate:"omitempty,datetime=2006-01-02" field:"date"`
	CategoryID int64   `json:"category_id" field:"categoryid"`
}
