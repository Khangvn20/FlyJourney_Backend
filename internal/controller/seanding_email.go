package controller

import (
    "net/http"
    "strconv"

    "github.com/Khangvn20/FlyJourney_Backend/internal/core/port/service"
    "github.com/gin-gonic/gin"
)

type EmailController struct {
    bookingEmailService service.BookingEmailService
}

func NewEmailController(bookingEmailService service.BookingEmailService) *EmailController {
    return &EmailController{
        bookingEmailService: bookingEmailService,
    }
}

func (ec *EmailController) SendBookingConfirmationEmailHandler(c *gin.Context) {
    bookingIDStr := c.Query("booking_id")
    if bookingIDStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "booking_id is required"})
        return
    }

    bookingID, err := strconv.ParseInt(bookingIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking_id"})
        return
    }

    response := ec.bookingEmailService.SendBookingConfirmationEmail(bookingID)
    if !response.Status {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": response.ErrorMessage,
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Booking confirmation email sent successfully",
    })
}