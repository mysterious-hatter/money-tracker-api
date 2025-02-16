package models

type Operation struct {
	ID 		   int64  `json:"id" field:"id"`
	Name       string `json:"name" field:"name"`
	WalletID   int64  `json:"wallet_id" field:"walletid"`
	Sum        int64  `json:"sum" field:"sum"`
	Date       string `json:"date" field:"date"`
	CategoryID int64  `json:"category_id" field:"categoryid"`
}