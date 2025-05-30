package util

import (
	"Authorization/internal/core/domain"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user *domain.User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   user.Email,
		"role":    user.Role,
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 2).Unix(),
	})

	key := os.Getenv("JWT_KEY")
	if key == "" {
		log.Fatal("JWT_KEY is not set")
	}
	jwtSecretKey := []byte(key)
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return ""
	}
	return tokenString
}

func ValidateToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpacted signing method")
		}
		key := os.Getenv("JWT_KEY")
		if key == "" {
			log.Fatal("JWT_KEY is not set")
		}
		jwtSecretKey := []byte(key)
		return jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, errors.New("could not parse claims")
}

func ValidateTokenFromRequest(r *http.Request) (jwt.MapClaims, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil, errors.New("no token")
	}

	return ValidateToken(cookie.Value)
}
