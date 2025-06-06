package service

import (
	"fmt"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/dto"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/port/service"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type tokenService struct {
	secretKey []byte
	expireDur time.Duration
}

// NewTokenService creates a new token service instance
func NewTokenService() service.TokenService {
	// Thời hạn token: 24 giờ
	expireDur := 24 * time.Hour
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	// Kiểm tra secret key có được cấu hình hay không
	if len(secretKey) == 0 {
		// Trong môi trường phát triển, có thể dùng key mặc định
		secretKey = []byte("fly_journey_secret_key")
	}
	return &tokenService{
		secretKey: secretKey,
		expireDur: expireDur,
	}
}

// GenerateToken tạo JWT token từ thông tin user
func (s *tokenService) GenerateToken(user *dto.User) (string, int64, error) {
	expireTime := time.Now().Add(s.expireDur)
	expireAt := expireTime.Unix()

	claims := jwt.MapClaims{
		"user_id": user.UserID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     expireAt,
	}

	// Tạo token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", 0, err
	}

	return tokenString, expireAt, nil
}

// ValidateToken xác thực token và trả về thông tin claims
func (s *tokenService) ValidateToken(tokenString string) (*dto.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra phương thức ký token
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Kiểm tra hết hạn
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, fmt.Errorf("token expired")
			}
		}

		// Lấy thông tin user từ claims
		userID, _ := claims["user_id"].(float64)
		email, _ := claims["email"].(string)
		role, _ := claims["role"].(string)

		return &dto.TokenClaims{
			UserID: int(userID),
			Email:  email,
			Role:   role,
		}, nil
	}

	return nil, fmt.Errorf("invalid token")
}
