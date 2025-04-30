package models

type User struct {
	Id       int64  `json:"id" field:"id"`
	Name     string `json:"name" validate:"max=50" field:"name"`
	Nickname string `json:"nickname" validate:"min=5,max=50" field:"nickname"`
	Password string `json:"password" validate:"min=5,max=50" field:"password"`
}
