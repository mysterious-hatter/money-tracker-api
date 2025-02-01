package storage

import "finances-backend/models"

type Storage interface {
	Open(host, username, passsword, dbname string) error
	Close()
	// User
	CreateUser(user *models.User) (int64, error)
	GetPasswordByEmail(email string) (string, error)
	GetAllUsers() (*[]models.User, error)
}
