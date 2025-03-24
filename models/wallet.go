package models

type Wallet struct {
	ID 	     int64  `json:"id" field:"id"`
	OwnerID  int64  `json:"owner_id" field:"ownerid"`
	Name     string `json:"name" validate:"min=5,max=50" field:"name"`
	Currency string `json:"currency" validate:"required,max=50" field:"currency"`
	Balance  int64  `json:"balance" field:"balance"`
}

// Сделать второй объект-валидатор, в котором будет проверка в зависмости от контекста употребления