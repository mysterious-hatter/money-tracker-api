package services

import (
	"errors"
	"finances-backend/models"
	"finances-backend/storage"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService struct {
	storage storage.Storage
}

func NewUserService(st storage.Storage) *UserService {
	us := UserService{st}
	return &us
}

func (us *UserService) CreateUser(user *models.User) (id int64, err error) {
	user.Password, err = hashPassword(user.Password)
	if err != nil {
		return 0, err
	}

	id, err = us.storage.CreateUser(user)
	return id, err
}

func (us *UserService) GetUserById(id int64) (*models.User, error) {
	user, err := us.storage.GetUserById(id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
