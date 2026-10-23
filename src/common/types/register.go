package types

import "time"

type RegisterRequest struct {
	Name            string `json:"name"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type RegisterResponse struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	UserId    Uuid      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func RegisterError(message string) RegisterResponse {
	return RegisterResponse{
		Success: false,
		Message: message,
	}
}
