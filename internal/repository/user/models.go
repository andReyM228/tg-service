package user

type loginRequest struct {
	ChatID   int64  `json:"chat_id"`
	Password string `json:"password"`
}
