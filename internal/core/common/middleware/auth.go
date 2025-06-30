package middleware

import (
    "log"
    "net/http"
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
        
        tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
        if tokenString == "" {
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "status": false,
                "errorCode": "UNAUTHORIZED",
                "errorMessage": "Token missing",
            })
            return
        }
       userID, role, err := tokenService.ValidateToken(tokenString)
        if err != nil {
            log.Printf(" Token validation failed: %v", err)
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "status": false,
                "errorCode": "UNAUTHORIZED",
                "errorMessage": "Invalid or expired token",
            })
            return
        }

        ctx.Set("userID", userID)
        ctx.Set("userRole", role)
        ctx.Next()
    }
}