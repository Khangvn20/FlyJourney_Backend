package middleware

import (
    "net/http"
    "strconv"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/port/service"
)

func AuthMiddleware(tokenService service.TokenService) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        authHeader := ctx.GetHeader("Authorization")
        if authHeader == "" {
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "status": false,
                "errorCode": "UNAUTHORIZED",
                "errorMessage": "Authorization header missing",
            })
            return
		}
        tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
        if tokenString == "" {
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "status": false,
                "errorCode": "UNAUTHORIZED",
                "errorMessage": "Token missing",
            })
            return
        }

        claims, err := tokenService.ValidateToken(tokenString)
        if err != nil {
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "status": false,
                "errorCode": "UNAUTHORIZED",
                "errorMessage": "Invalid or expired token",
            })
            return
        }

	     userIDStr, ok := claims["user_id"].(string)
        if !ok {
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "status": false,
                "errorCode": "UNAUTHORIZED",
                "errorMessage": "Invalid token payload",
            })
            return
        }
        userID, err := strconv.Atoi(userIDStr)
        if err != nil {
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "status": false,
                "errorCode": "UNAUTHORIZED",
                "errorMessage": "Invalid user ID in token",
            })
            return
        }

        ctx.Set("userID", userID)
        if role, exists := claims["role"].(string); exists {
            ctx.Set("userRole", role)
        }
        ctx.Next()
    }
}