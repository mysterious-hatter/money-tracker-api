package storage

import "finances-backend/models"

type Storage interface {
	Open(host, username, passsword, dbname string) error
	Close()
	// User
	GetAllUsers() (*[]models.Users, error)
}
