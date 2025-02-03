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
	// Wallet
	CreateWallet(wallet *models.Wallet) (int64, error)
}
