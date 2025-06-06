package service

import "github.com/Khangvn20/FlyJourney_Backend/internal/core/dto"

type TokenService interface {
	GenerateToken(user *dto.User) (string, int64, error)
	ValidateToken(token string) (*dto.TokenClaims, error)
}
