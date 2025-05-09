package models

type Operation struct {
	Id         int64    `json:"id" gorm:"<-:create"`
	Name       string   `json:"name" validate:"omitempty,max=50" gorm:"<-"`
	WalletId   int64    `json:"walletId" gorm:"<-:create;column:walletid"`
	Sum        float64  `json:"sum" validate:"nonzero" gorm:"<-"`
	Date       DateOnly `json:"date" gorm:"<-"`
	Place      string   `json:"place" validate:"omitempty,max=50" gorm:"<-"`
	CategoryId int64    `json:"categoryId" gorm:"<-;column:categoryid"`
}
