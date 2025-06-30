package router

import (
    "github.com/Khangvn20/FlyJourney_Backend/internal/controller"
    "github.com/gin-gonic/gin"
)
func AuthRoutes(rg *gin.RouterGroup, userController *controller.UserController,authMiddleware gin.HandlerFunc) {
	authRoutes := rg.Group("/auth")
	{
		authRoutes.POST("/register", userController.Register)
		authRoutes.POST("/login", userController.Login)
		authRoutes.POST("/confirm-register", userController.ConfirmRegister)
		authRoutes.POST("/confirm-reset-password", userController.ConfirmResetPassword)
		authRoutes.POST("/reset-password", userController.ResetPassword)
		authRoutes.POST("/logout", authMiddleware, userController.Logout)
 	}
}