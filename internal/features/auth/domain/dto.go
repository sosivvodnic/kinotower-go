package domain

type SignUpRequest struct {
	FIO      string `json:"fio"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Birthday string `json:"birthday"`
	GenderID int    `json:"gender_id"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
	ID     int    `json:"id"`
	FIO    string `json:"fio"`
}

type InvalidAuthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
