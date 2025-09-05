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
}