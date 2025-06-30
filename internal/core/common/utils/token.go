package utils

import (
    "log"
    "os"
    "errors"
    "time"
	"strconv"
    "github.com/golang-jwt/jwt/v5"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/port/service"
  
)

type jwtTokenService struct {
    secret []byte
    redisService service.RedisService
}

func NewTokenService(redisService service.RedisService) service.TokenService {
    secret := []byte(os.Getenv("JWT_SECRET"))
    return &jwtTokenService{
        secret:       secret,
        redisService: redisService,
    }
}

func (j *jwtTokenService) GenerateToken(userID int, role string, duration time.Duration) (string, error) {
    log.Printf("GenerateToken called with userID: %d, role: %s, duration: %v", userID, role, duration)
    
    if len(j.secret) == 0 {
        log.Printf("JWT secret is empty")
        return "", errors.New("JWT secret not configured")
    }
    
    claims := jwt.MapClaims{
        "user_id": strconv.Itoa(userID),
        "role":    role,      
        "exp":     time.Now().Add(duration).Unix(),
        "iat":     time.Now().Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(j.secret)
    if err != nil {
        return "", err
    }
     key := "token:" + tokenString
    err = j.redisService.Set(key, "valid", duration)
    if err !=nil {
        return "", err
    }
    return tokenString, nil
}
func (j *jwtTokenService) ValidateToken(tokenString string) (int, string, error) {
    revokedKey := "revoked:" + tokenString
    exists, err := j.redisService.Exists(revokedKey)
    if err != nil {
        return 0, "", err
    }
    if exists {
        return 0, "", jwt.ErrTokenInvalidClaims
    }
    
    tokenKey := "token:" + tokenString 
    exists, err = j.redisService.Exists(tokenKey)
    if err != nil {
        log.Printf("Error checking active token: %v", err) 
        return 0, "", err
    }
    if !exists {
        return 0, "", jwt.ErrTokenExpired
    }
    
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.ErrSignatureInvalid
        }
        return j.secret, nil
    })
    if err != nil {
        log.Printf("Error parsing token: %v", err)
        return 0, "", err
    }
    
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        userIDStr, ok := claims["user_id"].(string)
        if !ok {
            return 0, "", errors.New("invalid user_id in token")
        }
        
        userID, err := strconv.Atoi(userIDStr)
        if err != nil {
            return 0, "", errors.New("invalid user_id format")
        }
        
        role, ok := claims["role"].(string)
        if !ok {
            return 0, "", errors.New("invalid role in token")
        }
        
        return userID, role, nil
    }
    
    return 0, "", jwt.ErrSignatureInvalid
}
func (j *jwtTokenService) DeleteToken(tokenString string) error {
    
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return j.secret, nil
    })
    var ttl time.Duration = 24 * time.Hour
       if err == nil {
        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            if exp, ok := claims["exp"].(float64); ok {
                expTime := time.Unix(int64(exp), 0)
                ttl = time.Until(expTime)
                if ttl <= 0 {
                    ttl = time.Minute
                }
            }
        }
    }
    tokenKey := "token:" + tokenString 
    err = j.redisService.Del(tokenKey)
    if err != nil {
        log.Printf("Error deleting active token: %v", err)
        return err
    }
    log.Printf("Active token deleted successfully")
       revokedKey := "revoked:" + tokenString
    err = j.redisService.Set(revokedKey, "revoked", ttl) 
    if err != nil {
        log.Printf("Error adding to blacklist: %v", err)
        return err
    }
    log.Printf("Token added to blacklist successfully")
    
    return nil
}
