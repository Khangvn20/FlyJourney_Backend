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
	 CheckSeatOccupied(flightID int64, seatNumber string) (bool, error)
    CheckAlreadyCheckedIn(bookingDetailID int64) (bool, error)
    AssignSeatToBookingDetail(bookingDetailID int64, seatNumber string, flightID int64) error
    CreateCheckinRecord(checkinData *dto.CheckinDTO) (*dto.CheckinDTO, error)
    GetFlightInfoByBooking(bookingID int64) (*dto.FlightBasicInfo, error)
    ValidateBookingDetailOwnership(bookingDetailID int64, bookingID int64) (*dto.PassengerInfo, error)
}