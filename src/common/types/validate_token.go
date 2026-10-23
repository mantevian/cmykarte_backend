package types

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserId Uuid   `json:"user_id"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

type ValidateTokenRequest struct {
	Token string `json:"token"`
}

type ValidateTokenResponse struct {
	Success bool   `json:"success"`
	Claims  Claims `json:"claims"`
}

func ValidateTokenFalse() ValidateTokenResponse {
	return ValidateTokenResponse{Success: false}
}
