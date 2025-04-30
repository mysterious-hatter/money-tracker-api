package services

import (
	"errors"
	"finances-backend/models"
	"finances-backend/storage"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	storage       storage.Storage
	jwtSercet     string
	jwtExpiration int
}

func NewAuthService(st storage.Storage, jwtSec string, jwtExp int) *AuthService {
	as := AuthService{st, jwtSec, jwtExp}
	return &as
}

func (as *AuthService) AuthenticateUser(loginData *models.User) (string, error) {
	// Get user by Nickname
	user, err := as.storage.GetUserByNickname(loginData.Nickname)
	if err != nil {
		return "", err // user not found
	}
	// Check password
	if err := CheckPassword(loginData.Password, user.Password); err != nil {
		return "", err // wrong password
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"id":  user.Id,
		"exp": time.Now().Add(time.Hour * time.Duration(as.jwtExpiration)).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	encodedToken, err := token.SignedString([]byte(as.jwtSercet))

	return encodedToken, err
}

func (as *AuthService) AuthorizeUser(token interface{}) (int64, error) {
	// Decode token
	jwtToken, ok := token.(*jwt.Token)
	if !ok {
		return -1, errors.New("invalid token")
	}

	// Get payload
	payload, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return -1, errors.New("invalid token")
	}

	return int64(payload["id"].(float64)), nil
}

func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
