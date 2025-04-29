package models

type Operation struct {
	Id         int64   `json:"id" field:"id"`
	Name       string  `json:"name" validate:"omitempty,max=50" field:"name"`
	WalletId   int64   `json:"wallet_id" field:"walletid"`
	Sum        float64 `json:"sum" validate:"nonzero" field:"sum"`
	Date       string  `json:"date" validate:"omitempty,datetime=2006-01-02" field:"date"`
	Place      string  `json:"place" validate:"omitempty,max=50" field:"place"`
	CategoryId int64   `json:"category_id" field:"categoryid"`
}
