package storage

import (
	"finances-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresStorage struct {
	db *gorm.DB
}

func NewPostgresStorage() Storage {
	return &PostgresStorage{}
}

func (s *PostgresStorage) Open(host, username, passsword, dbname string) (err error) {
	dsn := "user=" + username + " password=" + passsword + " dbname=" + dbname + " host=" + host + " sslmode=disable"
	s.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return
}

// Users
func (s *PostgresStorage) CreateUser(user *models.User) (int64, error) {
	result := s.db.Table("users").Create(&user)

	return user.Id, result.Error
}

func (s *PostgresStorage) GetUserById(id int64) (*models.User, error) {
	user := models.User{}
	result := s.db.Table("users").First(&user, id)
	return &user, result.Error
}

func (s *PostgresStorage) GetUserByNickname(nickname string) (*models.User, error) {
	user := models.User{}
	result := s.db.Table("users").Where("nickname = ?", nickname).First(&user)
	return &user, result.Error
}

func (s *PostgresStorage) GetAllUsers() (*[]models.User, error) {
	users := []models.User{}
	result := s.db.Table("users").Find(&users)
	return &users, result.Error
}

// Wallets
func (s *PostgresStorage) CreateWallet(wallet *models.Wallet) (id int64, err error) {
	result := s.db.Table("wallets").Create(&wallet)
	return wallet.Id, result.Error
}

func (s *PostgresStorage) GetAllWallets(userId int64) ([]models.Wallet, error) {
	wallets := []models.Wallet{}
	result := s.db.Table("wallets").Select("wallets.*, COALESCE(SUM(operations.sum), 0) AS balance").
		Joins("LEFT JOIN operations ON wallets.id = operations.walletid").
		Where("wallets.ownerid = ?", userId).
		Group("wallets.id").
		Find(&wallets)
	return wallets, result.Error
}

func (s *PostgresStorage) GetWalletById(walletId int64) (*models.Wallet, error) {
	wallet := models.Wallet{}
	result := s.db.Table("wallets").Select("wallets.*, COALESCE(SUM(operations.sum), 0) AS balance").
		Joins("LEFT JOIN operations ON wallets.id = operations.walletid").
		Where("wallets.id = ?", walletId).
		Group("wallets.id").
		First(&wallet)
	return &wallet, result.Error
}

func (s *PostgresStorage) UpdateWallet(wallet *models.Wallet) error {
	result := s.db.Table("wallets").Save(wallet)
	return result.Error
}

// Categories
func (s *PostgresStorage) CreateCategory(category *models.Category) (id int64, err error) {
	result := s.db.Table("categories").Create(&category)
	return category.Id, result.Error
}

func (s *PostgresStorage) GetAllCategories(userId int64) ([]models.Category, error) {
	categories := []models.Category{}
	result := s.db.Table("categories").Where("ownerid = ?", userId).Find(&categories)
	return categories, result.Error
}

func (s *PostgresStorage) GetCategoryById(categoryId int64) (*models.Category, error) {
	category := models.Category{}
	result := s.db.Table("categories").First(&category, categoryId)
	return &category, result.Error
}

func (s *PostgresStorage) UpdateCategory(category *models.Category) error {
	result := s.db.Table("categories").Save(category)
	return result.Error
}

func (s *PostgresStorage) DeleteCategory(categoryId int64) error {
	result := s.db.Table("categories").Delete(&models.Category{}, categoryId)
	return result.Error
}

// Operations
func (s *PostgresStorage) CreateOperation(operation *models.Operation) (id int64, err error) {
	result := s.db.Table("operations").Create(&operation)
	return operation.Id, result.Error
}

func (s *PostgresStorage) GetOperations(walletId int64, sinceDate models.DateOnly, sortBy string) ([]models.Operation, error) {
	operations := []models.Operation{}

	basicQuery := s.db.Table("operations").Where("walletid = ?", walletId)
	if !sinceDate.IsZero() {
		basicQuery.Where("date >= ?", sinceDate)
	}
	if len(sortBy) > 0 {
		basicQuery.Order(sortBy + " DESC")
	}
	result := basicQuery.Find(&operations)

	return operations, result.Error
}

func (s *PostgresStorage) GetOperationById(operationId int64) (*models.Operation, error) {
	operation := models.Operation{}
	result := s.db.Table("operations").First(&operation, operationId)
	return &operation, result.Error
}

func (s *PostgresStorage) UpdateOperation(operation *models.Operation) error {
	updates := map[string]interface{}{}

	if operation.Name != "" {
		updates["name"] = operation.Name
	}
	if operation.Sum != 0.0 {
		updates["sum"] = operation.Sum
	}
	if !operation.Date.IsZero() {
		updates["date"] = operation.Date
	}
	if operation.Place != "" {
		updates["place"] = operation.Place
	}
	if operation.CategoryId != 0 {
		updates["categoryid"] = operation.CategoryId
	}

	if len(updates) == 0 {
		return nil
	}

	result := s.db.Table("operations").Where("id = ?", operation.Id).Updates(updates)
	return result.Error
}

func (s *PostgresStorage) DeleteOperation(operationId int64) error {
	result := s.db.Table("operations").Delete(&models.Operation{}, operationId)
	return result.Error
}
