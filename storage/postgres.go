package storage

import (
	"finances-backend/models"

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

func (s *PostgresStorage) Close() {
	s.db.Close()
}

// Users
func (s *PostgresStorage) CreateUser(user *models.User) (id int64, err error) {
	row := s.db.QueryRow("INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id",
		user.Name, user.Email, user.Password)
	err = row.Scan(&id)
	return
}

func (s *PostgresStorage) GetUserByID(id int64) (*models.User, error) {
	user := models.User{}
	res := s.db.QueryRowx("SELECT * FROM users WHERE id=$1", id)
	err := res.StructScan(&user)
	return &user, err
}

func (s *PostgresStorage) GetUserByEmail(email string) (*models.User, error) {
	user := models.User{}
	res := s.db.QueryRowx("SELECT * FROM users WHERE email=$1", email)
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
		wallet.Name, wallet.Currency, wallet.OwnerID)
	err = row.Scan(&id)
	return
}

func (s *PostgresStorage) GetAllWallets(userID int64) ([]models.Wallet, error) {
	wallets := []models.Wallet{}
	err := s.db.Select(&wallets,
		`SELECT wallets.*, COALESCE(SUM(operations.sum), 0) AS balance
		 FROM wallets LEFT JOIN operations
		 ON wallets.id=operations.walletid
		 WHERE wallets.ownerid=$1
		 GROUP BY wallets.id;`,
		userID)
	return wallets, err
}

func (s *PostgresStorage) GetWalletByID(walletID int64) (*models.Wallet, error) {
	wallet := models.Wallet{}
	res := s.db.QueryRowx(
		`SELECT wallets.*, COALESCE(SUM(operations.sum), 0) AS balance
		 FROM wallets
		 LEFT JOIN operations ON wallets.id = operations.walletid
         WHERE wallets.id = $1
         GROUP BY wallets.id;`,
		walletID)

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
		wallet.Name, wallet.Currency, wallet.ID)

	err := row.Scan(&wallet.Name, &wallet.Currency)
	return err
}

// Categories
func (s *PostgresStorage) CreateCategory(category *models.Category) (id int64, err error) {
	row := s.db.QueryRow("INSERT INTO categories (name, ownerid) VALUES ($1, $2) RETURNING id",
		category.Name, category.OwnerID)
	err = row.Scan(&id)
	return
}

func (s *PostgresStorage) GetAllCategories(userID int64) ([]models.Category, error) {
	categories := []models.Category{}
	err := s.db.Select(&categories, "SELECT * FROM categories WHERE ownerid=$1", userID)
	return categories, err
}

func (s *PostgresStorage) GetCategoryByID(categoryID int64) (*models.Category, error) {
	category := models.Category{}
	res := s.db.QueryRowx("SELECT * FROM categories WHERE id=$1", categoryID)
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
		category.Name, category.ID)

	err := row.Scan(&category.Name)
	return err
}

func (s *PostgresStorage) DeleteCategory(categoryID int64) error {
	_, err := s.db.Exec("DELETE FROM categories WHERE id=$1", categoryID)
	return err
}

// Operations
func (s *PostgresStorage) CreateOperation(operation *models.Operation) (id int64, err error) {
	row := s.db.QueryRow(
		`INSERT INTO operations (name, walletid, sum, date, place, categoryid) 
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		operation.Name, operation.WalletID, operation.Sum, operation.Date, operation.Place, operation.CategoryID)
	err = row.Scan(&id)
	return
}

func (s *PostgresStorage) GetOperationsByWalletID(walletID int64) ([]models.Operation, error) {
	operations := []models.Operation{}
	err := s.db.Select(&operations, "SELECT * FROM operations WHERE walletid=$1", walletID)
	return operations, err
}

func (s *PostgresStorage) GetOperationByID(operationID int64) (*models.Operation, error) {
	operation := models.Operation{}
	res := s.db.QueryRowx("SELECT * FROM operations WHERE id=$1", operationID)
	err := res.StructScan(&operation)
	return &operation, err
}

func (s *PostgresStorage) UpdateOperation(operation *models.Operation) error {
	row := s.db.QueryRow(
		`UPDATE operations
		 SET
		 	name = COALESCE(NULLIF($1, ''), name),
		 	sum = COALESCE(NULLIF($2, 0), sum),
		 	date = COALESCE(NULLIF($3, '')::DATE, date),
			place = COALESCE(NULLIF($4, ''), place),
		 	categoryid = COALESCE(NULLIF($5, 0), categoryid)
		 WHERE id=$6
		 RETURNING name, sum, date, place, categoryid;`,
		operation.Name, operation.Sum, operation.Date, operation.Place, operation.CategoryID, operation.ID)

	err := row.Scan(&operation.Name, &operation.Sum, &operation.Date, &operation.Place, &operation.CategoryID)
	return err
}

func (s *PostgresStorage) DeleteOperation(operationID int64) error {
	_, err := s.db.Exec("DELETE FROM operations WHERE id=$1", operationID)
	return err
}
