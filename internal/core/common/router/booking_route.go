package router

import (
	"github.com/Khangvn20/FlyJourney_Backend/internal/controller"
	"github.com/gin-gonic/gin"
)

func BookingRoutes(r *gin.RouterGroup, bookingController *controller.BookingController, authMiddleware gin.HandlerFunc) {
	bookingRoutes := r.Group("/booking")
	{
		bookingRoutes.POST("", authMiddleware, bookingController.CreateBooking)
		bookingRoutes.GET("/:bookingID", authMiddleware, bookingController.GetBookingByID)
		bookingRoutes.GET("/user/:userID", authMiddleware, bookingController.GetAllBookingByUserID)
	}
}
