package cars

type getCarByIDRequest struct {
	ID    int64  `json:"id"`
	Token string `json:"token"`
}
