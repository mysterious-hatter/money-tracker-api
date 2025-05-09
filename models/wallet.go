package models

type Wallet struct {
	Id       int64   `json:"id" gorm:"<-:create"`
	OwnerId  int64   `json:"ownerId" gorm:"<-:create;column:ownerid"`
	Name     string  `json:"name" validate:"omitempty,max=50" gorm:"<-"`
	Currency string  `json:"currency" validate:"omitempty,iso4217" gorm:"<-"`
	Balance  float64 `json:"balance" field:"balance"`
}
