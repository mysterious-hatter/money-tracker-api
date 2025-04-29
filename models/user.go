package models

type User struct {
	Id       int64  `json:"id" field:"id"`
	Name     string `json:"name" validate:"omitempty,min=5,max=50" field:"name"`
	Email    string `json:"email" validate:"email,max=50" field:"email"`
	Password string `json:"password" validate:"min=5,max=100" field:"password"`
}
