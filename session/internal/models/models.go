package models

type User struct {
	ID        int    `json:"id"`
	UUID      string `json:"uuid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Credentials struct {
	UserUUID string `json:"user_uuid"`
	Amount   int    `json:"amount"`
}

type Session struct {
	ID        int    `db:"id" json:"id"`
	BankID    string `db:"bank_id" json:"bank_id"`
	UserID    int    `db:"user_id" json:"user_id"`
	UserUUID  string `db:"user_uuid" json:"user_uuid"`
	Amount    int    `db:"amount" json:"amount"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
	Status    int    `db:"status" json:"status"` // 0 успешно //1 создана //2 в процессе у меня // 3 в процессе у банка // 4 ошибка
	Cancelled bool   `db:"cancelled" json:"cancelled"`
}

type UserBalance struct {
	UserID      int    `json:"user_id"`
	UserUUID    string `json:"user_uuid"`
	Balance     int    `json:"balance"`
	HoldBalance int    `json:"hold_balance"`
}

type DepositDTO struct {
	UserID   int    `json:"user_id"`
	UserUUID string `json:"user_uuid"`
	Amount   int    `json:"amount"`
}

type Operations struct {
	UserUUID string `json:"user_uuid"`
	BankUUID string `json:"bank_id"`
	Amount   int    `json:"amount"`
	Status   bool   `json:"status"`
}
