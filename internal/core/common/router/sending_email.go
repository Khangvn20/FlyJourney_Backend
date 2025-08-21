package router

import (
    "github.com/Khangvn20/FlyJourney_Backend/internal/controller"
    "github.com/gin-gonic/gin"
)

func BookingEmailRoute(api *gin.RouterGroup, emailController *controller.EmailController) {
    api.POST("/send-email", emailController.SendBookingConfirmationEmailHandler)
}