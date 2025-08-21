package service

import "github.com/Khangvn20/FlyJourney_Backend/internal/core/model/response"

type BookingEmailService interface {
    SendBookingConfirmationEmail(bookingID int64) *response.Response
}