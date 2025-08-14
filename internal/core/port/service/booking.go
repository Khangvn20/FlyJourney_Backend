package service

import (
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/model/request"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/model/response"
)

type BookingService interface {
	CreateBooking(req *request.CreateBookingRequest) *response.Response
	CancelExpiredBookings() *response.Response
}
