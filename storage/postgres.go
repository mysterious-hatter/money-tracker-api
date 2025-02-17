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
		 LEFT JOIN operations ON wallets.id=operations.walletid
         WHERE wallets.id=$1
         GROUP BY wallets.id;`,
		walletID)

	err := res.StructScan(&wallet)
	if err != nil {
		return nil, err
	}

	return &wallet, err
}
