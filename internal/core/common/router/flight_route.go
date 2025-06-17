package router

import (
    "github.com/Khangvn20/FlyJourney_Backend/internal/controller"
    "github.com/gin-gonic/gin"
)
func FlightRoutes(rg *gin.RouterGroup, flightController *controller.FlightController, authMiddleware gin.HandlerFunc) {
    flightRoutes := rg.Group("/flights")
    {
		flightRoutes.GET("/", flightController.GetAllFlights)
        flightRoutes.GET("/:id", flightController.GetFlightByID)
        flightRoutes.GET("/airline/:airline_id", flightController.GetFlightsByAirline)
        flightRoutes.GET("/status/:status", flightController.GetFlightsByStatus)
        flightRoutes.POST("/search", flightController.SearchFlights)
        //require authentication for the following routes
        flightRoutes.POST("", authMiddleware, flightController.CreateFlight)
        flightRoutes.PUT("/:id", authMiddleware, flightController.UpdateFlight)
        flightRoutes.PATCH("/:id/status", authMiddleware, flightController.UpdateFlightStatus)

    }
	}
	