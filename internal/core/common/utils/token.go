package utils

import (
    "os"
    "time"
	"strconv"
    "github.com/golang-jwt/jwt/v5"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/port/service"
)

type jwtTokenService struct {
    secret []byte
}

func NewTokenService() service.TokenService {
    secret := []byte(os.Getenv("JWT_SECRET"))
    return &jwtTokenService{secret: secret}
}

func (j *jwtTokenService) GenerateToken(userID int, duration time.Duration) (string, error) {
    claims := jwt.MapClaims{
        "user_id": strconv.Itoa(userID),
        "exp":     time.Now().Add(duration).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(j.secret)
}

func (j *jwtTokenService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.ErrSignatureInvalid
        }
        return j.secret, nil
    })
    if err != nil {
        return nil, err
    }
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims, nil
    }
    return nil, jwt.ErrSignatureInvalid
}