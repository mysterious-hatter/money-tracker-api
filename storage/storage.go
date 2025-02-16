package storage

import "finances-backend/models"

type Storage interface {
	Open(host, username, passsword, dbname string) error
	Close()
	// User
	CreateUser(user *models.User) (int64, error)
	GetUserByID(id int64) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetAllUsers() (*[]models.User, error)
	// UpdateUser(user *models.User) error
	// No delete user method because it is not needed

	// Wallet
	CreateWallet(wallet *models.Wallet) (int64, error)
	GetAllWallets(userID int64) ([]models.Wallet, error)
	GetWalletByID(walletID int64) (*models.Wallet, error)
	// UpdateWallet(wallet *models.Wallet) error
	// No delete wallet method because it is not needed

	// Transaction
	// CreateTransaction(transaction *models.Transaction) (int64, error)
	// GetTransactionsByWalletID(walletID int64) ([]models.Transaction, error)
}
