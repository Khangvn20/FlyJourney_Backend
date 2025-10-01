package service
import (
    "log"
    "strings"
	"fmt"
    "time"
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
     hasCheckedIn, err := s.checkinRepo.HasAnyCheckedInPassenger(pnrInfo.BookingID)
    if err != nil {
        log.Printf("Error checking if booking has checked in passengers: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Failed to verify check-in status",
        }
    }
    
    if hasCheckedIn {
        log.Printf("Booking %d has already checked in passengers", pnrInfo.BookingID)
        return &response.Response{
            Status:       true,
            ErrorCode:    error_code.Success,
            Data: &dto.CheckinValidationResponse{
                IsEligible:      false,
                ErrorMessage:    "This booking has already been checked in",
                PNRCode:         req.PNRCode,
                BookingID:       pnrInfo.BookingID,
                FlightID:        pnrInfo.FlightID,
                FlightNumber:    pnrInfo.FlightNumber,
                DepartureTime:   pnrInfo.DepartureTime,
                ArrivalTime:     pnrInfo.ArrivalTime,
                DepartureAirport: fmt.Sprintf("%s(%s)", pnrInfo.DepartureAirport, pnrInfo.DepartureCode),
                ArrivalAirport:   fmt.Sprintf("%s(%s)", pnrInfo.ArrivalAirport, pnrInfo.ArrivalCode),
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
        BookingID:        pnrInfo.BookingID,
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

func (s *checkinService) GetSeatMap(flightID int64) *response.Response {
    log.Printf("CheckSeats request for FlightID: %d", flightID)

    // Validate input
    if flightID <= 0 {
        log.Printf("Invalid FlightID: %d", flightID)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Invalid flight ID",
        }
    }


    seatCheckData, err := s.checkinRepo.GetConfirmedSeatsByFlightID(flightID)
    if err != nil {
        if err.Error() == "flight not found" {
            log.Printf("Flight not found for FlightID: %d", flightID)
            return &response.Response{
                Status:       false,
                ErrorCode:    error_code.NOTFOUND,
                ErrorMessage: "Flight not found",
            }
        }

        log.Printf("Error retrieving confirmed seats for FlightID %d: %v", flightID, err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Failed to retrieve confirmed seats",
        }
    }

    log.Printf("Confirmed seats retrieved successfully for FlightID: %d, Total seats: %d", 
        flightID, len(seatCheckData.ConfirmedSeats))

    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Seatmap data retrieved successfully",
        Data:         seatCheckData,
    }
}

func (s *checkinService) ProcessOnlineCheckin(req *request.OnlineCheckinRequest) *response.Response {
    log.Printf("ProcessOnlineCheckin request for BookingID: %d with %d passengers", 
        req.BookingID, len(req.Checkins))

    // 1. Validate basic request
    if req.BookingID <= 0 {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Invalid booking ID",
        }
    }

    if len(req.Checkins) == 0 {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "At least one passenger must be selected for check-in",
        }
    }

    // 2. Validate booking can be checked in
    flightInfo, validationError := s.validateBookingForCheckin(req.BookingID)
    if validationError != nil {
        return validationError
    }

    // 3. Process checkins using atomic repository method
    var allBoardingPasses []dto.BoardingPassInfo
    checkinTime := time.Now()

    for i, checkinReq := range req.Checkins {
        log.Printf("Processing passenger %d/%d - BookingDetailID: %d, SeatNumber: %s", 
            i+1, len(req.Checkins), checkinReq.BookingDetailID, checkinReq.SeatNumber)

        // Pre-validation checks
        if err := s.preValidateCheckin(checkinReq.BookingDetailID, req.BookingID, checkinReq.SeatNumber); err != nil {
            return err
        }

        // Generate boarding pass code and prepare atomic request
        boardingPassCode := s.generateBoardingPassCode()
        atomicReq := &dto.CheckinOnline{
            BookingDetailID:  checkinReq.BookingDetailID,
            SeatNumber:       checkinReq.SeatNumber,
            FlightID:         flightInfo.FlightID,
            BoardingPassCode: boardingPassCode,
            CheckinTime:      checkinTime,
        }

        // Execute atomic checkin operation
        result, err := s.checkinRepo.ProcessCheckinOnline(atomicReq)
        if err != nil {
            log.Printf("Error processing checkin for passenger %d (BookingDetailID %d): %v", 
                i+1, checkinReq.BookingDetailID, err)
            return s.handleCheckinError(err, checkinReq.SeatNumber)
        }

        // Collect boarding passes from result
        allBoardingPasses = append(allBoardingPasses, result.BoardingPasses...)
        
        log.Printf("Successfully processed passenger %d/%d", i+1, len(req.Checkins))
    }

    // 4. Build consolidated response
    return s.buildCheckinResponse(req.BookingID, flightInfo, allBoardingPasses, checkinTime)
}

// preValidateCheckin - Common validation before calling repository
func (s *checkinService) preValidateCheckin(bookingDetailID, bookingID int64, seatNumber string) *response.Response {
    // 1. Validate booking detail ownership
    passengerInfo, err := s.checkinRepo.ValidateBookingDetailOwnership(bookingDetailID, bookingID)
    if err != nil {
        log.Printf("Ownership validation failed for BookingDetailID %d: %v", bookingDetailID, err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: fmt.Sprintf("Booking detail %d does not belong to this booking", bookingDetailID),
        }
    }

    // 2. Check if already checked in
    alreadyCheckedIn, err := s.checkinRepo.CheckAlreadyCheckedIn(bookingDetailID)
    if err != nil {
        log.Printf("Error checking if already checked in: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Failed to validate check-in status",
        }
    }
    
    if alreadyCheckedIn {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: fmt.Sprintf("Passenger %s has already checked in", passengerInfo.PassengerName),
        }
    }

    // 3. Basic seat validation
    if seatNumber == "" {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Seat number is required",
        }
    }

    return nil // No errors
}

// handleCheckinError - Centralized error handling
func (s *checkinService) handleCheckinError(err error, seatNumber string) *response.Response {
    errorMsg := err.Error()
    
    if strings.Contains(errorMsg, "already occupied") {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: fmt.Sprintf("Seat %s is already occupied", seatNumber),
        }
    }
    
    if strings.Contains(errorMsg, "not found") {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.NOTFOUND,
            ErrorMessage: "Booking or passenger information not found",
        }
    }
    
    if strings.Contains(errorMsg, "transaction") {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Database transaction error. Please try again.",
        }
    }
    
    return &response.Response{
        Status:       false,
        ErrorCode:    error_code.InternalError,
        ErrorMessage: "Failed to process check-in. Please try again.",
    }
}


// buildCheckinResponse - Xây dựng response cho check-in thành công
func (s *checkinService) buildCheckinResponse(bookingID int64, flightInfo *dto.FlightBasicInfo, boardingPasses []dto.BoardingPassInfo, checkinTime time.Time) *response.Response {
    // Calculate boarding time
    boardingTime := flightInfo.DepartureTime.Add(-30 * time.Minute)

    // Build consolidated response
    finalresponse := &dto.OnlineCheckinResponse{
        BookingID:      bookingID,
        FlightNumber:   flightInfo.FlightNumber,
        CheckinTime:    checkinTime,
        BoardingTime:   &boardingTime,
        CheckedInCount: len(boardingPasses),
        BoardingPasses: boardingPasses,
    }

    log.Printf("Successfully completed checkin for %d passengers in booking %d", len(boardingPasses), bookingID)

    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: fmt.Sprintf("Successfully checked in %d passengers", len(boardingPasses)),
        Data:         finalresponse,
    }
}

func (s *checkinService) validateBookingForCheckin(bookingID int64) (*dto.FlightBasicInfo, *response.Response) {
    log.Printf("Validating booking for checkin - BookingID: %d", bookingID)

    // Get flight info by booking ID
    flightInfo, err := s.checkinRepo.GetFlightInfoByBooking(bookingID)
    if err != nil {
        if strings.Contains(err.Error(), "not found") {
            log.Printf("Booking not found: %d", bookingID)
            return nil, &response.Response{
                Status:       false,
                ErrorCode:    error_code.NOTFOUND,
                ErrorMessage: "Booking not found",
            }
        }
        log.Printf("Error getting flight info for BookingID %d: %v", bookingID, err)
        return nil, &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Failed to retrieve booking information",
        }
    }

    // Validate booking status
    if flightInfo.BookingStatus != "confirmed" && flightInfo.BookingStatus != "paid" {
        log.Printf("Invalid booking status: %s for BookingID: %d", flightInfo.BookingStatus, bookingID)
        return nil, &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: fmt.Sprintf("Booking status must be confirmed or paid, current: %s", flightInfo.BookingStatus),
        }
    }

    // Validate checkin time window (optional - uncomment if needed)
    /*
    now := time.Now()
    checkinOpenTime := flightInfo.DepartureTime.Add(-24 * time.Hour)
    checkinCloseTime := flightInfo.DepartureTime.Add(-2 * time.Hour)

    if now.Before(checkinOpenTime) {
        return nil, &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: fmt.Sprintf("Check-in opens 24 hours before departure at %s", 
                checkinOpenTime.Format("2006-01-02 15:04")),
        }
    }

    if now.After(checkinCloseTime) {
        return nil, &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: fmt.Sprintf("Check-in closed 2 hours before departure at %s", 
                checkinCloseTime.Format("2006-01-02 15:04")),
        }
    }
    */

    log.Printf("Booking validation successful for BookingID: %d, FlightID: %d", bookingID, flightInfo.FlightID)
    return flightInfo, nil
}
// generateBoardingPassCode - Tạo mã boarding pass
func (s *checkinService) generateBoardingPassCode() string {
    return fmt.Sprintf("BP%d%d", time.Now().Unix()%10000, time.Now().Nanosecond()%1000)
}