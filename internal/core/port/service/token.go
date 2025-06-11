package service

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
    GenerateToken(userID int, duration time.Duration) (string, error)
    ValidateToken(tokenString string) (jwt.MapClaims, error)
}