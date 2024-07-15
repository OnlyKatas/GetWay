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
	ID     int    `json:"id"`
	UUID   string `json:"uuid"`
	Amount int    `json:"amount"`
}

type Session struct {
	ID        int    `json:"id"`
	UUID      string `json:"uuid"`
	UserID    int    `json:"user_id"`
	UserUUID  string `json:"user_uuid"`
	Amount    int    `json:"amount"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Status    int    `json:"status"` // 0 успешно //1 создана //2 в процессе у меня // 3 в процессе у банка // 4 ошибка
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
