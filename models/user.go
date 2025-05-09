package models

type User struct {
	Id       int64  `json:"id" gorm:"<-:create"`
	Name     string `json:"name" validate:"max=50" gorm:"<-"`
	Nickname string `json:"nickname" validate:"min=5,max=50" gorm:"<-"`
	Password string `json:"password" validate:"min=5,max=50" gorm:"<-"`
}
