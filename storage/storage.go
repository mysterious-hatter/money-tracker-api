package storage

import (
	"finances-backend/models"
	"time"
)

type Storage interface {
	Open(host, username, passsword, dbname string) error
	Close() error
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
	UpdateWallet(wallet *models.Wallet) error
	// No delete wallet method because it is not needed

	// Category
	CreateCategory(category *models.Category) (int64, error)
	GetAllCategories(userID int64) ([]models.Category, error)
	GetCategoryByID(categoryID int64) (*models.Category, error)
	UpdateCategory(category *models.Category) error
	DeleteCategory(categoryID int64) error
	
	// Operation
	CreateOperation(operation *models.Operation) (int64, error)
	GetOperationsByWalletID(walletID int64) ([]models.Operation, error)
	GetOperationsSinceDateByWalletID(walletID int64, date time.Time) ([]models.Operation, error)
	GetOperationByID(operationID int64) (*models.Operation, error)
	UpdateOperation(operation *models.Operation) error
	DeleteOperation(operationID int64) error
}
