package service

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
    "context"
)

type TokenService interface {
    GenerateToken(userID int, duration time.Duration) (string, error)
    ValidateToken(tokenString string) (jwt.MapClaims, error)
    DeleteToken(tokenString string) error
}
type RevokedTokenRepository interface {
    RevokeToken(ctx context.Context, token string, userID int, expiryAt time.Time) error
    IsTokenRevoked(ctx context.Context, token string) (bool, error)
    CleanupExpiredTokens(ctx context.Context) (int, error)
    RevokeAllUserTokens(ctx context.Context, userID int) error
}