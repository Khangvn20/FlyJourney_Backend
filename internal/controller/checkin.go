package controller 

import (
    "net/http"
    "log"
    "fmt"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/entity/error_code"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/model/request"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/port/service"
)

type CheckinController struct {
    checkinService service.CheckinService
}

func NewCheckinController(checkinService service.CheckinService) *CheckinController {
    return &CheckinController{
        checkinService: checkinService,
    }
}

// ValidateCheckin - API validate checkin by PNR
func (c *CheckinController) ValidateCheckin(ctx *gin.Context) {
    log.Printf("ValidateCheckin API called")

    var req request.ValidateCheckin
    if err := ctx.ShouldBindJSON(&req); err != nil {
        log.Printf("Invalid request body: %v", err)
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":       false,
            "errorCode":    error_code.InvalidRequest,
            "errorMessage": "Invalid request format: " + err.Error(),
        })
        return
    }

    // Validate required fields
    if req.PNRCode == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":       false,
            "errorCode":    error_code.InvalidRequest,
            "errorMessage": "PNR code is required",
        })
        return
    }

    if req.Email == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":       false,
            "errorCode":    error_code.InvalidRequest,
            "errorMessage": "Email is required",
        })
        return
    }

    if req.FullName == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":       false,
            "errorCode":    error_code.InvalidRequest,
            "errorMessage": "Full name is required",
        })
        return
    }

    log.Printf("Validating checkin for PNR: %s, Email: %s", req.PNRCode, req.Email)

    // Call service
    result := c.checkinService.ValidateCheckin(&req)

    // Determine HTTP status code based on result
    statusCode := http.StatusOK
    if !result.Status {
        switch result.ErrorCode {
        case error_code.InvalidRequest:
            statusCode = http.StatusBadRequest
        case error_code.NOTFOUND:
            statusCode = http.StatusNotFound
        case error_code.InternalError:
            statusCode = http.StatusInternalServerError
        default:
            statusCode = http.StatusInternalServerError
        }
    }

    ctx.JSON(statusCode, result)
}

// CheckSeats - API check confirmed seats by flight ID

func (c *CheckinController) GetSeatMap(ctx *gin.Context) {
    flightIDStr := ctx.Param("flight_id")
    flightID, err := strconv.ParseInt(flightIDStr, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":       false,
            "errorCode":    error_code.InvalidRequest,
            "errorMessage": "Invalid flight ID format",
        })
        return
    }

    result := c.checkinService.GetSeatMap(flightID)

    statusCode := http.StatusOK
    if !result.Status {
        switch result.ErrorCode {
        case error_code.InvalidRequest:
            statusCode = http.StatusBadRequest
        case error_code.NOTFOUND:
            statusCode = http.StatusNotFound
        case error_code.InternalError:
            statusCode = http.StatusInternalServerError
        default:
            statusCode = http.StatusInternalServerError
        }
    }

    ctx.JSON(statusCode, result)
}

func (c *CheckinController) ProcessOnlineCheckin(ctx *gin.Context) {
        log.Printf("ProcessOnlineCheckin API called")

    // 1. Parse request body
    var req request.OnlineCheckinRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        log.Printf("Invalid request body: %v", err)
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":       false,
            "errorCode":    error_code.InvalidRequest,
            "errorMessage": "Invalid request format: " + err.Error(),
        })
        return
    }

    // 2. Validate required fields
    if req.BookingID <= 0 {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":       false,
            "errorCode":    error_code.InvalidRequest,
            "errorMessage": "Booking ID is required",
        })
        return
    }

    if len(req.Checkins) == 0 {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":       false,
            "errorCode":    error_code.InvalidRequest,
            "errorMessage": "At least one check-in detail is required",
        })
        return
    }

    // Validate each check-in detail
    for i, checkin := range req.Checkins {
        if checkin.BookingDetailID <= 0 {
            ctx.JSON(http.StatusBadRequest, gin.H{
                "status":       false,
                "errorCode":    error_code.InvalidRequest,
                "errorMessage": fmt.Sprintf("Invalid booking detail ID at index %d", i),
            })
            return
        }

        if checkin.SeatNumber == "" {
            ctx.JSON(http.StatusBadRequest, gin.H{
                "status":       false,
                "errorCode":    error_code.InvalidRequest,
                "errorMessage": fmt.Sprintf("Seat number is required at index %d", i),
            })
            return
        }
    }

    log.Printf("Processing online check-in for BookingID: %d with %d passengers", req.BookingID, len(req.Checkins))

    // 3. Call service
    result := c.checkinService.ProcessOnlineCheckin(&req)

    // 4. Determine HTTP status code based on result
    statusCode := http.StatusOK
    if !result.Status {
        switch result.ErrorCode {
        case error_code.InvalidRequest:
            statusCode = http.StatusBadRequest
        case error_code.NOTFOUND:
            statusCode = http.StatusNotFound
        case error_code.InternalError:
            statusCode = http.StatusInternalServerError
        default:
            statusCode = http.StatusInternalServerError
        }
    }

    // 5. Return response
    ctx.JSON(statusCode, result)
}