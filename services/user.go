package services

import (
	"finances-backend/models"
	"finances-backend/storage"

	"golang.org/x/crypto/bcrypt"
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

func (us *UserService) AuthorizeUser(user *models.User) error {
	hashedPassword, err := us.storage.GetPasswordByEmail(user.Email)
	if err != nil {
		return err // user not found
	}

	return CheckPassword(user.Password, hashedPassword)
}

func (us *UserService) GetAllUsers() (*[]models.User, error) {
	users, err := us.storage.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func HashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}

func CheckPassword(password string, hashedPassword string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}