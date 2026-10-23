package types

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	UserId  Uuid   `json:"user_id"`
	Token   string `json:"token"`
}

func LoginError(message string) LoginResponse {
	return LoginResponse{
		Success: false,
		Message: message,
	}
}
