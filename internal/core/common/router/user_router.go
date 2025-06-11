package router

import (
	"github.com/Khangvn20/FlyJourney_Backend/internal/controller"
	"github.com/gin-gonic/gin"
)

func SetupRouter(userController *controller.UserController) *gin.Engine {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// API versioning và grouping
	api := r.Group("/api")
	v1 := api.Group("/v1")

	// User routes
	userRoutes := v1.Group("/users")
	{
		userRoutes.POST("/register", userController.Register)
		userRoutes.POST("/login", userController.Login)
		userRoutes.POST("/confirm-register", userController.ConfirmRegister)
		userRoutes.POST("/reset-password", userController.ResetPassword)
		userRoutes.POST("/confirm-reset-password", userController.ConfirmResetPassword)
	}

	return r
}
