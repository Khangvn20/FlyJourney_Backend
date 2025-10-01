package repository
import (
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/dto"
	
)
type CheckinRepository interface {
	//Validate checkin
	 GetPNRInfo(pnrCode string) (*dto.PNRInfo, error)
	 GetBookingDetailsForCheckin(bookingID int64) ([]*dto.CheckinBookingDetail, error)
	//GetSeatmap
	GetConfirmedSeatsByFlightID(flightID int64) (*dto.SeatCheckResponse, error)
	//Checkin
    ProcessCheckinOnline(checkinRequest *dto.CheckinOnline) (*dto.OnlineCheckinResponse, error)
    ValidateBookingDetailOwnership(bookingDetailID int64, bookingID int64) (*dto.PassengerInfo, error)
	 CheckAlreadyCheckedIn(bookingDetailID int64) (bool, error)  
    GetFlightInfoByBooking(bookingID int64) (*dto.FlightBasicInfo, error)
	HasAnyCheckedInPassenger(bookingID int64) (bool, error)
}