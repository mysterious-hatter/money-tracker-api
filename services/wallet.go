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
		return 0, ErrUnableToCreate
	}
	return id, err
}

func (ws *WalletService) GetAllWallets(userId int64) ([]models.Wallet, error) {
	wallets, err := ws.storage.GetAllWallets(userId)
	if err != nil {
		return nil, ErrNoWalletsFound
	}
	return wallets, err
}

func (ws *WalletService) GetWalletById(walletId int64, userId int64) (*models.Wallet, error) {
	wallet, err := ws.storage.GetWalletById(walletId)
	if err != nil {
		return nil, ErrWalletNotFound
	}

	// Check user's ownership of the wallet
	if err := checkOwnership(wallet.OwnerId, userId); err != nil {
		return nil, err
	}
	return wallet, err
}

func (ws *WalletService) UpdateWallet(wallet *models.Wallet, userId int64) error {
	// Check user's ownership of the wallet and if it exists
	_, err := ws.GetWalletById(wallet.Id, userId)
	if err != nil {
		return err
	}

	err = ws.storage.UpdateWallet(wallet)
	if err != nil {
		return ErrUnableToUpdate
	}
	return err
}

// Deletion of wallets is not forseen,
// as it would lead to deletion all the connected operations.
