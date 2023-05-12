package user

type loginRequest struct {
	ChatID   int64  `json:"chat_id"`
	Password string `json:"password"`
}

type loginResponse struct {
	UserID int64 `json:"user_id"`
}
