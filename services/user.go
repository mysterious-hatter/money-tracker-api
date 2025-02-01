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
func (us *UserService) GetAllUsers() (*[]models.Users, error) {
	users, err := us.storage.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}
