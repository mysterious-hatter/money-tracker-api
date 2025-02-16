package services

import (
	"errors"
	"finances-backend/models"
	"finances-backend/storage"
)

var (
	ErrWalletNotFound        error = errors.New("wallet not found")
	ErrWalletNotBelongToUser error = errors.New("requested wallet does not belong to the user")
	ErrNoWalletsFound        error = errors.New("no wallets found")
	ErrSomethingWentWrong    error = errors.New("something went wrong")
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
	if err != nil {
		return 0, ErrSomethingWentWrong
	}
	return id, err
}

func (ws *WalletService) GetAllWallets(userID int64) ([]models.Wallet, error) {
	wallets, err := ws.storage.GetAllWallets(userID)
	if err != nil {
		return nil, ErrNoWalletsFound
	}
	return wallets, err
}

func (ws *WalletService) GetWalletByID(walletID int64, userID int64) (*models.Wallet, error) {
	wallet, err := ws.storage.GetWalletByID(walletID)
	if err != nil {
		return nil, ErrWalletNotFound
	}

	if wallet.OwnerID != userID {
		return nil, ErrWalletNotBelongToUser
	}
	return wallet, err
}