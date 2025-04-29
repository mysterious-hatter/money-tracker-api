package models

type Category struct {
	ID      int64  `json:"id" field:"id"`
	Name    string `json:"name" field:"name"`
	OwnerID int64  `json:"owner_id" field:"ownerid"`
}
