package storage

import (
	"finances-backend/models"
)

type Storage interface {
	Open(host, username, passsword, dbname string) error
	// Close() error
	// User
	CreateUser(user *models.User) (int64, error)
	GetUserById(id int64) (*models.User, error)
	GetUserByNickname(nickname string) (*models.User, error)
	GetAllUsers() (*[]models.User, error)
	// UpdateUser(user *models.User) error
	// No delete user method because it is not needed

	// Wallet
	CreateWallet(wallet *models.Wallet) (int64, error)
	GetAllWallets(userId int64) ([]models.Wallet, error)
	GetWalletById(walletId int64) (*models.Wallet, error)
	UpdateWallet(wallet *models.Wallet) error
	// No delete wallet method because it is not needed

	// Category
	CreateCategory(category *models.Category) (int64, error)
	GetAllCategories(userId int64) ([]models.Category, error)
	GetCategoryById(categoryId int64) (*models.Category, error)
	UpdateCategory(category *models.Category) error
	DeleteCategory(categoryId int64) error

	// Operation
	CreateOperation(operation *models.Operation) (int64, error)
	GetOperations(walletId int64, sinceDate models.DateOnly, sortBy string) ([]models.Operation, error)
	GetOperationById(operationId int64) (*models.Operation, error)
	UpdateOperation(operation *models.Operation) error
	DeleteOperation(operationId int64) error

	// Search
	SearchOperations(name string, walletId int64, date models.DateOnly,
		place string, categoryId int64, sortBy string) ([]models.Operation, error)
}
