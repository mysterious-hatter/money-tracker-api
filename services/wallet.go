package services

import (
	"finances-backend/models"
	"finances-backend/storage"
)

type WalletService struct {
	storage storage.Storage
}

func NewWalletService(st storage.Storage) *WalletService {
	ws := WalletService{st}
	return &ws
}

func (ws *WalletService) CreateWallet(wallet *models.Wallet) (int64, error) {
	id, err := ws.storage.CreateWallet(wallet)
	return id, err
}