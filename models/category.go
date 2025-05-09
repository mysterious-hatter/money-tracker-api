package models

type Category struct {
	Id      int64  `json:"id" gorm:"<-:create"`
	Name    string `json:"name" validate:"min=1,max=50" gorm:"<-"`
	OwnerId int64  `json:"ownerId" gorm:"<-:create;column:ownerid"`
}
