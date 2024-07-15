package models

type Request struct {
	UserUUID string `json:"user_uuid" db:"user_uuid"`
	BankUUID string `json:"bank_uuid" db:"bank_uuid"`
	Amount   int    `json:"amount" db:"amount"`
	Status   int    `json:"status" db:"status"` // 0 in progress, 1 successfully, 2 error
}
