package router

import (
    "github.com/Khangvn20/FlyJourney_Backend/internal/controller"
    "github.com/gin-gonic/gin"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/common/middleware"
)
func FlightRoutes(rg *gin.RouterGroup, flightController *controller.FlightController, authMiddleware gin.HandlerFunc) {
    flightRoutes := rg.Group("/flights")
    {
		flightRoutes.GET("/", flightController.GetAllFlights)
        flightRoutes.GET("/:id", flightController.GetFlightByIDForUser)
        flightRoutes.GET("/airline/:airline_id", flightController.GetFlightsByAirline)
        flightRoutes.GET("/status/:status", flightController.GetFlightsByStatus)
        flightRoutes.POST("/search", flightController.SearchFlightsForUser)
         flightRoutes.POST("/search/roundtrip", flightController.SearchRoundtripFlightsForUser)
    }
    adminRoutes :=rg.Group("/admin/flights")
    adminRoutes.Use(authMiddleware, middleware.RequireAdmin())
    {
        adminRoutes.POST("/", flightController.CreateFlight)
        adminRoutes.POST("/:id/classes", flightController.CreateFlightClasses)
        adminRoutes.PUT("/:id", flightController.UpdateFlight)
        adminRoutes.PATCH("/:id", flightController.UpdateFlightStatus)
        adminRoutes.GET("/:id", flightController.GetFlightByID)
        adminRoutes.POST("/search", flightController.SearchFlights)
        adminRoutes.POST("/search/roundtrip", flightController.SearchRoundtripFlights)
    }
	}
	