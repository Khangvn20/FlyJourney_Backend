package controller
import (
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/model/request"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/port/service"
    "github.com/gin-gonic/gin"
    "net/http"
    "strconv"
)
type FlightController struct {
	flightService service.FlightService
}
func NewFlightController(flightService service.FlightService) *FlightController {
    return &FlightController{
        flightService: flightService,
    }
}
func (c *FlightController) CreateFlight(ctx *gin.Context) {
    var req request.CreateFlightRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":       false,
            "errorCode":    "INVALID_REQUEST",
            "errorMessage": err.Error(),
        })
        return
    }

    result := c.flightService.CreateFlight(&req)
    
    statusCode := http.StatusOK
    if !result.Status {
        statusCode = http.StatusBadRequest
    }
    
    ctx.JSON(statusCode, result)
}
func (c *FlightController) GetFlightByID(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":       false,
            "errorCode":    "INVALID_ID",
            "errorMessage": "Invalid flight ID format",
        })
        return
    }

    result := c.flightService.GetFlightByID(id)
    
    statusCode := http.StatusOK
    if !result.Status {
        statusCode = http.StatusNotFound
    }
    
    ctx.JSON(statusCode, result)
}
func (c *FlightController) GetAllFlights(ctx *gin.Context) {
    pageStr := ctx.DefaultQuery("page", "1")
    limitStr := ctx.DefaultQuery("limit", "10")
    
    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        page = 1
    }
    
    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit < 1 {
        limit = 10
    }
    
    result := c.flightService.GetAllFlights(page, limit)
    ctx.JSON(http.StatusOK, result)
}
func (c *FlightController) UpdateFlight(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":       false,
            "errorCode":    "INVALID_ID",
            "errorMessage": "Invalid flight ID format",
        })
        return
    }
    
    var req request.UpdateFlightRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":       false,
            "errorCode":    "INVALID_REQUEST",
            "errorMessage": err.Error(),
        })
        return
    }
    
    result := c.flightService.UpdateFlight(id, &req)
    
    statusCode := http.StatusOK
    if !result.Status {
        statusCode = http.StatusBadRequest
        if result.ErrorCode == "NOT_FOUND" {
            statusCode = http.StatusNotFound
        }
    }
    
    ctx.JSON(statusCode, result)
}
func (c *FlightController) SearchFlights(ctx *gin.Context) {
    var req request.FlightSearchRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":       false,
            "errorCode":    "INVALID_REQUEST",
            "errorMessage": err.Error(),
        })
        return
    }
    
    result := c.flightService.SearchFlights(&req)
    ctx.JSON(http.StatusOK, result)
}
func (c *FlightController) GetFlightsByAirline(ctx *gin.Context) {
    airlineIDStr := ctx.Param("airline_id")
    airlineID, err := strconv.Atoi(airlineIDStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":       false,
            "errorCode":    "INVALID_AIRLINE_ID",
            "errorMessage": "Invalid airline ID format",
        })
        return
    }
    
    pageStr := ctx.DefaultQuery("page", "1")
    limitStr := ctx.DefaultQuery("limit", "10")
    
    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        page = 1
    }
    
    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit < 1 {
        limit = 10
    }
    
    result := c.flightService.GetFlightByAirline(airlineID, page, limit)
    ctx.JSON(http.StatusOK, result)
}
func (c *FlightController) GetFlightsByStatus(ctx *gin.Context) {
    status := ctx.Param("status")
    
    pageStr := ctx.DefaultQuery("page", "1")
    limitStr := ctx.DefaultQuery("limit", "10")
    
    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        page = 1
    }
    
    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit < 1 {
        limit = 10
    }
    
    result := c.flightService.GetFlightsByStatus(status, page, limit)
    ctx.JSON(http.StatusOK, result)
}
func (c *FlightController) UpdateFlightStatus(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":       false,
            "errorCode":    "INVALID_ID",
            "errorMessage": "Invalid flight ID format",
        })
        return
    }
    
    var req request.UpdateFlightStatusRequest
    
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":       false,
            "errorCode":    "INVALID_REQUEST",
            "errorMessage": err.Error(),
        })
        return
    }
    
    result := c.flightService.UpdateFlightStatus(id, &req)
    
    statusCode := http.StatusOK
    if !result.Status {
        statusCode = http.StatusBadRequest
        if result.ErrorCode == "NOT_FOUND" {
            statusCode = http.StatusNotFound
        }
    }
    
    ctx.JSON(statusCode, result)
}
func (c *FlightController) SearchRoundtripFlights(ctx *gin.Context) {
    var req request.RoundtripFlightSearchRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":       false,
            "errorCode":    "INVALID_REQUEST",
            "errorMessage": err.Error(),
        })
        return
    }
    
    result := c.flightService.SearchRoundtripFlights(&req)
    ctx.JSON(http.StatusOK, result)
}