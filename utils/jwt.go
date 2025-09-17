package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "supersecret"

func GenerateToken(userId int64, email, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    userId,
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(time.Hour * 2).Unix(),
	})
	return token.SignedString([]byte(secretKey))
}

func VerifyToken(tokenString string) (int64, string, string, error) {
	parsedToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(secretKey), nil
	})
	if err != nil || !parsedToken.Valid {
		return 0, "", "", errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", "", errors.New("invalid token")
	}

	id := int64(claims["id"].(float64))
	email := claims["email"].(string)
	role := claims["role"].(string)

	return id, email, role, nil
}
