package controller

import (
	"net/http"
    "strconv"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/entity/error_code"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/model/request"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/port/service"
	"github.com/gin-gonic/gin"
)

type BookingController struct {
    bookingService service.BookingService
}

func NewBookingController(bookingService service.BookingService) *BookingController {
    return &BookingController{
        bookingService: bookingService,
    }
}

func (c *BookingController) CreateBooking(ctx *gin.Context) {
     userID, exists := ctx.Get("userID")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{
            "status":       false,
            "errorCode":    error_code.InternalErrMsg,
            "errorMessage": "Không tìm thấy thông tin người dùng",
        })
        return
    }
	var req request.CreateBookingRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":       false,
            "errorCode":    error_code.InvalidRequest,
            "errorMessage": err.Error(),
        })
        return
    }
	 req.UserID = int64(userID.(int)) 

    result := c.bookingService.CreateBooking(&req)
    
    statusCode := http.StatusCreated
    if !result.Status {
        switch result.ErrorCode {
        case error_code.ResourceLocked:
            statusCode = http.StatusConflict // 409 Conflict
        case error_code.NoAvailableSeats, error_code.InsufficientSeats, error_code.InvalidRequest:
            statusCode = http.StatusBadRequest // 400 Bad Request
        case error_code.NOTFOUND:
            statusCode = http.StatusNotFound // 404 Not Found
        default:
            statusCode = http.StatusInternalServerError // 500 Internal Server Error
        }
    }

    ctx.JSON(statusCode, result)
}

func (c *BookingController) GetBookingByID(ctx *gin.Context) {
    bookingIDParam := ctx.Param("bookingID")
    bookingID, err := strconv.ParseInt(bookingIDParam, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":  false,
            "message": "Invalid booking ID",
        })
        return
    }

    response := c.bookingService.GetBookingID(bookingID)
    if !response.Status {
        ctx.JSON(http.StatusNotFound, gin.H{
            "status":  false,
            "message": response.ErrorMessage,
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "status":  true,
        "message": response.ErrorMessage,
        "data":    response.Data,
    })
}