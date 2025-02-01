package storage

import (
	"finances-backend/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sqlx.DB
}

func NewPostgesStorage() Storage {
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

func (s *PostgresStorage) GetAllUsers() (*[]models.Users, error) {
	return nil, nil
}