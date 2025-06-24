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
    revokedTokens  map[string]struct{}
}

func NewTokenService() service.TokenService {
    secret := []byte(os.Getenv("JWT_SECRET"))
    return &jwtTokenService{secret: secret , revokedTokens: make(map[string]struct{})}
}

func (j *jwtTokenService) GenerateToken(userID int, role string, duration time.Duration) (string, error) {
    claims := jwt.MapClaims{
        "user_id": strconv.Itoa(userID),
        "role":    role,      
        "exp":     time.Now().Add(duration).Unix(),
        "iat":     time.Now().Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(j.secret)
}
func (j *jwtTokenService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
      if _, revoked := j.revokedTokens[tokenString]; revoked {
        return nil, jwt.ErrSignatureInvalid 
    }
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
func (j *jwtTokenService) DeleteToken(tokenString string) error {
    j.revokedTokens[tokenString] = struct{}{}
    return nil
}