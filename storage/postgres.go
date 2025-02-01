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

func (s *PostgresStorage) CreateUser(user *models.User) (id int64, err error) {
	row := s.db.QueryRow("INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id",
		user.Name, user.Email, user.Password)
	err = row.Scan(&id)
	return
}

func (s *PostgresStorage) GetPasswordByEmail(email string) (string, error) {
	var password string
	err := s.db.Get(&password, "SELECT password FROM users WHERE email=$1", email)
	return password, err
}

func (s *PostgresStorage) GetAllUsers() (*[]models.User, error) {
	users := []models.User{}
	err := s.db.Select(&users, "SELECT * FROM users")
	return &users, err
}