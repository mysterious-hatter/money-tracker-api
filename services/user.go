package services

import (
	"finances-backend/models"
	"finances-backend/storage"
)

type UserService struct {
	storage storage.Storage
}

func NewUserService(st storage.Storage) *UserService {
	us := UserService{st}
	return &us
}

func (us *UserService) GetUserByID(id int64) (*models.User, error) {
	user, err := us.storage.GetUserByID(id)
	return user, err
}
