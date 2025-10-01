package controller

import (
	"net/http"
	"strconv"
    "log"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/model/response"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/entity/error_code"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/model/request"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/port/service"
	"github.com/gin-gonic/gin"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/common/utils"
)
type FlightController struct {
	flightService service.FlightService
    bookingEmailService service.BookingEmailService
}
func NewFlightController(flightService service.FlightService, bookingEmailService service.BookingEmailService) *FlightController {
    return &FlightController{
        flightService: flightService,
        bookingEmailService: bookingEmailService,
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
func (c *FlightController) BatchCreateFlights(ctx *gin.Context) {
    var req request.BatchCreateFlightRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":       false,
            "errorCode":    "INVALID_REQUEST",
            "errorMessage": err.Error(),
        })
        return
    }

    result := c.flightService.BatchCreateFlights(&req)

    statusCode := http.StatusOK
    if !result.Status {
        statusCode = http.StatusBadRequest
    }

    ctx.JSON(statusCode, result)
}
func (c *FlightController) CreateFlightClasses(ctx *gin.Context) {
     flightIDStr := ctx.Param("id")
    flightID, err := strconv.Atoi(flightIDStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":       false,
            "errorCode":    "INVALID_ID",
            "errorMessage": "Invalid flight ID format",
        })
        return
    }
    var req []request.FlightClassRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":       false,
            "errorCode":    "INVALID_REQUEST",
            "errorMessage": err.Error(),
        })
        return
    }

    result, err := c.flightService.CreateFlightClasses(flightID, req)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "status":       false,
            "errorCode":    "INTERNAL_ERROR", 
            "errorMessage": "Failed to create flight classes",
        })
        return
    }
    
    statusCode := http.StatusCreated
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
      userRole, roleExists := ctx.Get("userRole")
     var result *response.Response
    
      if roleExists && userRole == "{admin}" {
        result = c.flightService.GetFlightByIDForAdmin(id)
    } else {
        result = c.flightService.GetFlightByIDForUser(id)
    }
    
    statusCode := http.StatusOK
    if !result.Status {
        if result.ErrorCode == error_code.InternalError {
            statusCode = http.StatusNotFound
        } else {
            statusCode = http.StatusInternalServerError
        }
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
func (c *FlightController) GetFareCLassCode(ctx *gin.Context) {
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
    
    result := c.flightService.GetFareCLassCode(id)
    
    statusCode := http.StatusOK
    if !result.Status {
        statusCode = http.StatusNotFound
    }
    
    ctx.JSON(statusCode, result)
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
func (c *FlightController) GetFlightByIDForUser(ctx *gin.Context) {
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
    result := c.flightService.GetFlightByIDForUser(id)
    statusCode := http.StatusOK
    if !result.Status {
        statusCode = http.StatusBadRequest
        if result.ErrorCode == "NOT_FOUND" {
            statusCode = http.StatusNotFound
        }
    }
    
    ctx.JSON(statusCode, result)
}
func (c *FlightController) GetFlightsByDate(ctx *gin.Context) {
    var req request.GetFlightsByDateRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":       false,
            "errorCode":    "INVALID_REQUEST",
            "errorMessage": err.Error(),
        })
        return
    }
    
    result := c.flightService.GetFlightsByDate(&req)
    ctx.JSON(http.StatusOK, result)
}
func (c *FlightController) SearchFlightsForUser(ctx *gin.Context) {
    var req request.FlightSearchRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":       false,
            "errorCode":    "INVALID_REQUEST",
            "errorMessage": err.Error(),
        })
        return
    }
    
    result := c.flightService.SearchFlightsForUser(&req)
    ctx.JSON(http.StatusOK, result)
}
func (c *FlightController) SearchRoundtripFlightsForUser(ctx *gin.Context) {
    var req request.RoundtripFlightSearchRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":       false,
            "errorCode":    "INVALID_REQUEST",
            "errorMessage": err.Error(),
        })
        return
    }
    
    result := c.flightService.SearchRoundtripFlightsForUser(&req)
    ctx.JSON(http.StatusOK, result)
}


func (c *FlightController) UpdateFlightTime(ctx *gin.Context) {
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

    var req request.UpdateFlightTimeRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":       false,
            "errorCode":    "INVALID_REQUEST",
            "errorMessage": err.Error(),
        })
        return
    }

    result := c.flightService.UpdateFlightTime(int64(id), &req)
    ctx.JSON(http.StatusOK, result)
}

func (c *FlightController) QueueDelayNotifications (ctx *gin.Context){
    var req struct {
        FlightID        int64  `json:"flight_id"`
        NewDepartureTime string `json:"new_departure_time"`
        Reason           string `json:"reason"`
    }
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Invalid request format: " + err.Error(),
        })
        return
    }

 
    newDepartureTime, err := utils.ParseTime(req.NewDepartureTime)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Invalid time format. Expected format: dd/MM/yyyy HH:mm",
        })
        return
    }
    if req.FlightID <= 0 {
        ctx.JSON(http.StatusBadRequest, response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Invalid flight ID",
        })
        return
    }

    if len(req.Reason) < 5 {
        ctx.JSON(http.StatusBadRequest, response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Reason must be at least 5 characters",
        })
        return
    }

    log.Printf("Queuing delay notifications for flight %d with new departure time %s and reason: %s",
        req.FlightID, req.NewDepartureTime, req.Reason)


    result := c.bookingEmailService.QueueFlightDelayNotifications(req.FlightID, newDepartureTime, req.Reason)

    ctx.JSON(http.StatusOK, result)
}