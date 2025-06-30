package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/entity/error_code"
)

func RequireAdmin() gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("userRole")
        if !exists {
            c.JSON(http.StatusForbidden, gin.H{
                "status":       false,
                "errorCode":    error_code.Forbidden,
                "errorMessage": "Permission denied: user role not found",
            })
            c.Abort()
            return
        }
        roleStr:= userRole.(string)
        if roleStr != "{admin}" {
            c.JSON(http.StatusForbidden, gin.H{
                "status":       false,
                "errorCode":    error_code.Forbidden,
                "errorMessage": "This action is only allowed for administrators",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}