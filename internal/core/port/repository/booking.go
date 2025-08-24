package repository
import (
	
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/dto"
)
type BookingRepository interface {
    // Create booking
    CreateBooking(booking *dto.Booking) (*dto.Booking, error)
    CancelBookings(bookingIDs []int64) ([]int64, error)
    // Get bookings
    GetExpiredBookingIDs() ([]int64, error) 
    GetBookingByID(bookingID int64) (*dto.Booking, error)
    UpdateStatusConfirm(bookingID int64) (*dto.Booking, error)
    CheckFlightClassAvailability(flightClassID int64) (bool, int, error)
    UpdateBookingStatus(bookingID int64, status string) (*dto.Booking, error)
    GetBookingsByFlightID(flightID int64) ([]*dto.Booking, error)
    GetAllBookingByUserID(userID int64) ([]*dto.Booking, error)
}