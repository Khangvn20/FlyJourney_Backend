package service
import (
    "log"
    "strings"
	"fmt"
//	"time"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/dto"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/entity/error_code"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/model/request"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/model/response"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/port/repository"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/port/service"
)
type checkinService struct {
	checkinRepo repository.CheckinRepository
	bookingRepo repository.BookingRepository
}

func NewCheckinService(checkinRepo repository.CheckinRepository, bookingRepo repository.BookingRepository) service.CheckinService {
	return &checkinService{
		checkinRepo: checkinRepo,
		bookingRepo: bookingRepo,
	}
}
func (s *checkinService) ValidateCheckin(req *request.ValidateCheckin) *response.Response {
    log.Printf("ValidateCheckin request: PNR=%s, Email=%s, Name=%s", 
        req.PNRCode, req.Email, req.FullName)

    // 1. Lấy PNR info từ repository (data only)
    pnrInfo, err := s.checkinRepo.GetPNRInfo(req.PNRCode)
    if err != nil {
        if err.Error() == "PNR not found" {
            
            return &response.Response{
                Status:       true,
                ErrorCode:    error_code.Success,
                Data: &dto.CheckinValidationResponse{
                    IsEligible:   false,
                    ErrorMessage: "PNR code not found",
                    PNRCode:      req.PNRCode,
                },
            }
        }
        
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Failed to retrieve PNR information",
        }
    }

    // 2. BUSINESS VALIDATION - Email match
    if !strings.EqualFold(pnrInfo.ContactEmail, req.Email) {
        return &response.Response{
            Status:       true,
            ErrorCode:    error_code.Success,
            Data: &dto.CheckinValidationResponse{
                IsEligible:   false,
                ErrorMessage: "Email does not match booking record",
                PNRCode:      req.PNRCode,
            },
        }
    }

    // 3. BUSINESS VALIDATION - Name match
    if !strings.Contains(strings.ToLower(pnrInfo.ContactName), strings.ToLower(req.FullName)) {
        return &response.Response{
            Status:       true,
            ErrorCode:    error_code.Success,
            Data: &dto.CheckinValidationResponse{
                IsEligible:   false,
                ErrorMessage: "Name does not match booking record",
                PNRCode:      req.PNRCode,
            },
        }
    }

    // 4. BUSINESS VALIDATION - PNR status
    if pnrInfo.Status != "active" {
        return &response.Response{
            Status:       true,
            ErrorCode:    error_code.Success,
            Data: &dto.CheckinValidationResponse{
                IsEligible:   false,
                ErrorMessage: "PNR is not active",
                PNRCode:      req.PNRCode,
            },
        }
    }

    // 5. BUSINESS VALIDATION - Booking status
    if pnrInfo.BookingStatus != "confirmed" && pnrInfo.BookingStatus != "paid" {
        return &response.Response{
            Status:       true,
            ErrorCode:    error_code.Success,
            Data: &dto.CheckinValidationResponse{
                IsEligible:   false,
                ErrorMessage: fmt.Sprintf("Booking must be confirmed or paid, current: %s", pnrInfo.BookingStatus),
                PNRCode:      req.PNRCode,
                FlightNumber: pnrInfo.FlightNumber,
            },
        }
    }

    // 6. BUSINESS VALIDATION - Flight status
    if pnrInfo.FlightStatus == "cancelled" {
        return &response.Response{
            Status:       true,
            ErrorCode:    error_code.Success,
            Data: &dto.CheckinValidationResponse{
                IsEligible:   false,
                ErrorMessage: "Flight has been cancelled",
                PNRCode:      req.PNRCode,
                FlightNumber: pnrInfo.FlightNumber,
            },
        }
    }

    // 7. BUSINESS VALIDATION - Time window
/*    now := time.Now()
    checkinOpenTime := pnrInfo.DepartureTime.Add(-24 * time.Hour)
    checkinCloseTime := pnrInfo.DepartureTime.Add(-2 * time.Hour)

   if now.Before(checkinOpenTime) {
        return &response.Response{
            Status:       true,
            ErrorCode:    error_code.Success,
            Data: &dto.CheckinValidationResponse{
                IsEligible:   false,
                ErrorMessage: fmt.Sprintf("Checkin opens 24 hours before departure at %s", 
                    checkinOpenTime.Format("2006-01-02 15:04")),
                PNRCode:      req.PNRCode,
                FlightNumber: pnrInfo.FlightNumber,
                DepartureTime: pnrInfo.DepartureTime,
                ArrivalTime:  pnrInfo.ArrivalTime,
            },
        }
    }

    if now.After(checkinCloseTime) {
        return &response.Response{
            Status:       true,
            ErrorCode:    error_code.Success,
            Data: &dto.CheckinValidationResponse{
                IsEligible:   false,
                ErrorMessage: fmt.Sprintf("Checkin closed 2 hours before departure at %s", 
                    checkinCloseTime.Format("2006-01-02 15:04")),
                PNRCode:      req.PNRCode,
                FlightNumber: pnrInfo.FlightNumber,
                DepartureTime: pnrInfo.DepartureTime,
                ArrivalTime:  pnrInfo.ArrivalTime,
                 DepartureAirport: fmt.Sprintf("%s(%s)", pnrInfo.DepartureAirport, pnrInfo.DepartureCode),
                ArrivalAirport:   fmt.Sprintf("%s(%s)", pnrInfo.ArrivalAirport, pnrInfo.ArrivalCode),
                AirlineName:      pnrInfo.AirlineName,
            },
        }
    }
*/
    // 8. Lấy booking details
    bookingDetails, err := s.checkinRepo.GetBookingDetailsForCheckin(pnrInfo.BookingID)
    if err != nil {
        log.Printf("Error retrieving booking details for BookingID %d: %v", pnrInfo.BookingID, err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Failed to retrieve passenger information",
        }
    }

    if len(bookingDetails) == 0 {
        return &response.Response{
            Status:       true,
            ErrorCode:    error_code.Success,
            Data: &dto.CheckinValidationResponse{
                IsEligible:   false,
                ErrorMessage: "No passengers found for this booking",
                PNRCode:      req.PNRCode,
                FlightNumber: pnrInfo.FlightNumber,
            },
        }
    }

    // 9. SUCCESS - All validations passed
    checkinResponse := &dto.CheckinValidationResponse{
        IsEligible:       true,
        PNRCode:          req.PNRCode,
        FlightID:         pnrInfo.FlightID,
        FlightNumber:     pnrInfo.FlightNumber,
        DepartureTime:    pnrInfo.DepartureTime,
        ArrivalTime:      pnrInfo.ArrivalTime,
        DepartureAirport: fmt.Sprintf("%s(%s)", pnrInfo.DepartureAirport, pnrInfo.DepartureCode),
        ArrivalAirport:   fmt.Sprintf("%s(%s)", pnrInfo.ArrivalAirport, pnrInfo.ArrivalCode),
        AirlineName:      pnrInfo.AirlineName,
      //  CheckinOpenTime:  &checkinOpenTime,
     //   CheckinCloseTime: &checkinCloseTime,
        BookingDetails:   bookingDetails,
    }

    log.Printf("Checkin validation successful for PNR: %s, Flight: %s, Passengers: %d", 
        req.PNRCode, pnrInfo.FlightNumber, len(bookingDetails))

    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Checkin validation successful",
        Data:         checkinResponse,
    }
}
