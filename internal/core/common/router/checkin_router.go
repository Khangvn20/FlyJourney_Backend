package router

import (
    "github.com/Khangvn20/FlyJourney_Backend/internal/controller"
    "github.com/gin-gonic/gin"
)
func CheckinRoutes(r *gin.RouterGroup, checkinController *controller.CheckinController) {
    checkinRoutes := r.Group("/checkin")
    {
        checkinRoutes.POST("/validate", checkinController.ValidateCheckin)
	}  
} 