package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"mantevian.xyz/cmykarte_backend/common/types"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(userId types.Uuid, name string) (string, error) {
	claims := types.Claims{
		UserId: userId,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*types.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &types.Claims{}, func(t *jwt.Token) (any, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*types.Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
