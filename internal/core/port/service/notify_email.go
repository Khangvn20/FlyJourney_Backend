package service

import (
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/model/response"
    "time"
)

type BookingEmailService interface {
    SendBookingConfirmationEmail(bookingID int64) *response.Response
    SendFlightCancelEmail(bookingID int64, reason string) *response.Response
    SendFlightDelayEmail(bookingID int64, newDepartureTime time.Time, reason string) *response.Response
}