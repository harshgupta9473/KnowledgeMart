package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type JWTClaims struct {
	ID   uint   `json:"userId"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func CreateJWT(userID uint, role string) (string, error) {
	secretKey := os.Getenv("JWTSECRETKEY")

	claims := &JWTClaims{
		ID:   userID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJWT(tokenString string) (*JWTClaims, error) {
	secretKey := os.Getenv("JWTSECRETKEY")
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		if claims.ExpiresAt.Before(time.Now()) {
			return nil, errors.New("token expired")
		}
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
