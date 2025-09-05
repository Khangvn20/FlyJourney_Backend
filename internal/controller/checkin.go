package controller

import (
    "net/http"
    "log"

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