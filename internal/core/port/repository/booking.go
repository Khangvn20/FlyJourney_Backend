package repository
import (
	
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/dto"
)
type BookingRepository interface {
    // Create booking
    CreateBooking(booking *dto.Booking) (*dto.Booking, error)
    
    // Get bookings
  //  GetBookingByID(bookingID int64) (*dto.Booking, error)
    //GetBookingsByUserID(userID int64, page, limit int) ([]*dto.Booking, int, error)
   // GetPendingPaymentBookings(timeThreshold time.Time) ([]int64, error)
    
    // Update booking
  /*  UpdateBookingStatus(bookingID int64, status string) error
    CancelBooking(bookingID int64) error
    ConfirmBooking(bookingID int64) error
    
    // Booking details
    UpdateBookingDetails(details []*dto.BookingDetail) error
    AssignSeat(bookingDetailID int64, seatID int64) error
    
    // Payment related
    UpdateBookingAfterPayment(bookingID int64, paymentID int64, status string) error
    
    
    // Statistics
    CountBookingsByStatus(status string) (int, error)
    GetRecentBookings(limit int) ([]*dto.Booking, error)*/
	// Check availability
	CheckFlightClassAvailability(flightClassID int64) (bool, int, error)
}