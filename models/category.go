package models

type Category struct {
	Id      int64  `json:"id" field:"id"`
	Name    string `json:"name" field:"name"`
	OwnerId int64  `json:"owner_id" field:"ownerid"`
}
