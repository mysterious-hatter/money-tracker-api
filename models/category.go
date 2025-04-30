package models

type Category struct {
	Id      int64  `json:"id" field:"id"`
	Name    string `json:"name" validate:"min=1,max=50" field:"name"`
	OwnerId int64  `json:"owner_id" field:"ownerid"`
}
