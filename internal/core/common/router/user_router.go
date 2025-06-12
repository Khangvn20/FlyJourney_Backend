package router

import (
    "github.com/Khangvn20/FlyJourney_Backend/internal/controller"
    "github.com/gin-gonic/gin"
)

func UserRoutes(rg *gin.RouterGroup, userController *controller.UserController,authMiddleware gin.HandlerFunc) {
    userRoutes := rg.Group("/users")
    {
        userRoutes.POST("/register", userController.Register)
        userRoutes.POST("/login", userController.Login)
        userRoutes.POST("/confirm-register", userController.ConfirmRegister)
        userRoutes.POST("/reset-password", userController.ResetPassword)
        userRoutes.POST("/confirm-reset-password", userController.ConfirmResetPassword)
        userRoutes.POST("/logout", userController.Logout)
        userRoutes.GET("/", authMiddleware, userController.GetUserInfo)
    }
}