package storage

import (
	"finances-backend/models"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sqlx.DB
}

func NewPostgresStorage() Storage {
	return &PostgresStorage{}
}

func (s *PostgresStorage) Open(host, username, passsword, dbname string) (err error) {
	connStr := "user=" + username + " password=" + passsword + " dbname=" + dbname + " host=" + host + " sslmode=disable"
	s.db, err = sqlx.Connect("postgres", connStr)
	return
}

func (s *PostgresStorage) Close() error {
	return s.db.Close()
}

// Users
func (s *PostgresStorage) CreateUser(user *models.User) (id int64, err error) {
	row := s.db.QueryRow("INSERT INTO users (name, nickname, password) VALUES ($1, $2, $3) RETURNING id",
		user.Name, user.Nickname, user.Password)
	err = row.Scan(&id)
	return
}

func (s *PostgresStorage) GetUserById(id int64) (*models.User, error) {
	user := models.User{}
	res := s.db.QueryRowx("SELECT * FROM users WHERE id=$1", id)
	err := res.StructScan(&user)
	return &user, err
}

func (s *PostgresStorage) GetUserByNickname(nickname string) (*models.User, error) {
	user := models.User{}
	res := s.db.QueryRowx("SELECT * FROM users WHERE nickname=$1", nickname)
	err := res.StructScan(&user)
	return &user, err
}

func (s *PostgresStorage) GetAllUsers() (*[]models.User, error) {
	users := []models.User{}
	err := s.db.Select(&users, "SELECT * FROM users")
	return &users, err
}

// Wallets
func (s *PostgresStorage) CreateWallet(wallet *models.Wallet) (id int64, err error) {
	row := s.db.QueryRow("INSERT INTO wallets (name, currency, ownerid) VALUES ($1, $2, $3) RETURNING id",
		wallet.Name, wallet.Currency, wallet.OwnerId)
	err = row.Scan(&id)
	return
}

func (s *PostgresStorage) GetAllWallets(userId int64) ([]models.Wallet, error) {
	wallets := []models.Wallet{}
	err := s.db.Select(&wallets,
		`SELECT wallets.*, COALESCE(SUM(operations.sum), 0) AS balance
		 FROM wallets LEFT JOIN operations
		 ON wallets.id=operations.walletid
		 WHERE wallets.ownerid=$1
		 GROUP BY wallets.id;`,
		userId)
	return wallets, err
}

func (s *PostgresStorage) GetWalletById(walletId int64) (*models.Wallet, error) {
	wallet := models.Wallet{}
	res := s.db.QueryRowx(
		`SELECT wallets.*, COALESCE(SUM(operations.sum), 0) AS balance
		 FROM wallets
		 LEFT JOIN operations ON wallets.id = operations.walletid
         WHERE wallets.id = $1
         GROUP BY wallets.id;`,
		walletId)

	err := res.StructScan(&wallet)
	if err != nil {
		return nil, err
	}

	return &wallet, err
}

func (s *PostgresStorage) UpdateWallet(wallet *models.Wallet) error {
	row := s.db.QueryRow(
		`UPDATE wallets
		 SET
		 	name = COALESCE(NULLIF($1, ''), name), 
		 	currency = COALESCE(NULLIF($2, ''), currency)
		 WHERE id = $3
		 RETURNING name, currency;`,
		wallet.Name, wallet.Currency, wallet.Id)

	err := row.Scan(&wallet.Name, &wallet.Currency)
	return err
}

// Categories
func (s *PostgresStorage) CreateCategory(category *models.Category) (id int64, err error) {
	row := s.db.QueryRow("INSERT INTO categories (name, ownerid) VALUES ($1, $2) RETURNING id",
		category.Name, category.OwnerId)
	err = row.Scan(&id)
	return
}

func (s *PostgresStorage) GetAllCategories(userId int64) ([]models.Category, error) {
	categories := []models.Category{}
	err := s.db.Select(&categories, "SELECT * FROM categories WHERE ownerid=$1", userId)
	return categories, err
}

func (s *PostgresStorage) GetCategoryById(categoryId int64) (*models.Category, error) {
	category := models.Category{}
	res := s.db.QueryRowx("SELECT * FROM categories WHERE id=$1", categoryId)
	err := res.StructScan(&category)
	return &category, err
}

func (s *PostgresStorage) UpdateCategory(category *models.Category) error {
	row := s.db.QueryRow(
		`UPDATE categories
		 SET
		 	name = COALESCE(NULLIF($1, ''), name)
		 WHERE id=$2
		 RETURNING name;`,
		category.Name, category.Id)

	err := row.Scan(&category.Name)
	return err
}

func (s *PostgresStorage) DeleteCategory(categoryId int64) error {
	_, err := s.db.Exec("DELETE FROM categories WHERE id=$1", categoryId)
	return err
}

// Operations
func (s *PostgresStorage) CreateOperation(operation *models.Operation) (id int64, err error) {
	row := s.db.QueryRow(
		`INSERT INTO operations (name, walletid, sum, date, place, categoryid) 
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		operation.Name, operation.WalletId, operation.Sum, operation.Date, operation.Place, operation.CategoryId)
	err = row.Scan(&id)
	return
}

func (s *PostgresStorage) GetOperationsByWalletId(walletId int64) ([]models.Operation, error) {
	operations := []models.Operation{}
	err := s.db.Select(&operations, "SELECT * FROM operations WHERE walletid=$1 ORDER BY Date DESC", walletId)
	return operations, err
}

func (s *PostgresStorage) GetOperationsSinceDateByWalletId(walletId int64, date time.Time) ([]models.Operation, error) {
	operations := []models.Operation{}
	err := s.db.Select(&operations,
		`SELECT * FROM operations WHERE walletid=$1 AND date >= $2 ORDER BY Date DESC`,
		walletId, date)
	return operations, err
}

func (s *PostgresStorage) GetOperationById(operationId int64) (*models.Operation, error) {
	operation := models.Operation{}
	res := s.db.QueryRowx("SELECT * FROM operations WHERE id=$1", operationId)
	err := res.StructScan(&operation)
	return &operation, err
}

func (s *PostgresStorage) UpdateOperation(operation *models.Operation) error {
	row := s.db.QueryRow(
		`UPDATE operations
		 SET
		 	name = COALESCE(NULLIF($1, ''), name),
		 	sum = COALESCE(NULLIF($2::FLOAT, 0.0), sum),
		 	date = COALESCE(NULLIF($3, '')::DATE, date),
			place = COALESCE(NULLIF($4, ''), place),
		 	categoryid = COALESCE(NULLIF($5, 0), categoryid)
		 WHERE id=$6
		 RETURNING name, sum, date, place, categoryid;`,
		operation.Name, operation.Sum, operation.Date, operation.Place, operation.CategoryId, operation.Id)

	err := row.Scan(&operation.Name, &operation.Sum, &operation.Date, &operation.Place, &operation.CategoryId)
	return err
}

func (s *PostgresStorage) DeleteOperation(operationId int64) error {
	_, err := s.db.Exec("DELETE FROM operations WHERE id=$1", operationId)
	return err
}
