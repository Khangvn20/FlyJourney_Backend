package router

import (
    "github.com/Khangvn20/FlyJourney_Backend/internal/controller"
    "github.com/gin-gonic/gin"
)

func UserRoutes(rg *gin.RouterGroup, userController *controller.UserController,authMiddleware gin.HandlerFunc) {
    userRoutes := rg.Group("/users")
    {
        userRoutes.GET("/", authMiddleware, userController.GetUserInfo)
        userRoutes.PUT("/", authMiddleware, userController.UpdateProfile)
    }
}