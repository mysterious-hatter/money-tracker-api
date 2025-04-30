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

func (us *UserService) CreateUser(user *models.User) (id int64, err error) {
	user.Password, err = HashPassword(user.Password)
	if err != nil {
		return 0, err
	}

	id, err = us.storage.CreateUser(user)
	return id, err
}

func (us *UserService) GetUserById(id int64) (*models.User, error) {
	user, err := us.storage.GetUserById(id)
	return user, err
}
