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

// ProcessOnlineCheckin - Xử lý check-in online
func (s *checkinService) ProcessOnlineCheckin(req *request.OnlineCheckinRequest) *response.Response {
    log.Printf("ProcessOnlineCheckin request for BookingID: %d with %d passengers", 
        req.BookingID, len(req.Checkins))

    // 1. Validate booking có thể check-in
    flightInfo, validationError := s.validateBookingForCheckin(req.BookingID)
    if validationError != nil {
        return validationError
    }

    // 2. Process each checkin request
    boardingPasses, processError := s.processCheckinRequests(req, flightInfo)
    if processError != nil {
        return processError
    }

    // 3. Build success response
    return s.buildCheckinResponse(req.BookingID, flightInfo, boardingPasses)
}

func (s *checkinService) validateBookingForCheckin(bookingID int64) (*dto.FlightBasicInfo, *response.Response) {    
     
    flightInfo, err := s.checkinRepo.GetFlightInfoByBooking(bookingID)
    if err != nil {
        if err.Error() == "booking not found" {
            log.Printf("Booking not found: %d", bookingID)
            return nil, &response.Response{
                Status:       false,
                ErrorCode:    error_code.NOTFOUND,
                ErrorMessage: "Booking not found",
            }
        }
        log.Printf("Error getting flight info: %v", err)
        return nil, &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Failed to retrieve booking information",
        }
    }
    if flightInfo.BookingStatus != "confirmed" && flightInfo.BookingStatus != "paid" {
        log.Printf("Invalid booking status: %s for BookingID: %d", flightInfo.BookingStatus, bookingID)
        return nil, &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: fmt.Sprintf("Booking status must be confirmed or paid, current: %s", flightInfo.BookingStatus),
        }
    }
  //  now := time.Now()
   // checkinCloseTime := flightInfo.DepartureTime.Add(-2 * time.Hour)
    
    /*if now.After(checkinCloseTime) {
        log.Printf("Check-in closed for BookingID: %d, current: %s, close time: %s", 
            bookingID, now.Format("2006-01-02 15:04"), checkinCloseTime.Format("2006-01-02 15:04"))
        return nil, &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: fmt.Sprintf("Check-in closed 2 hours before departure at %s", 
                checkinCloseTime.Format("2006-01-02 15:04")),
        }
    }*/

    log.Printf("Booking validation successful for BookingID: %d", bookingID)
    return flightInfo, nil
}

func (s *checkinService) processCheckinRequests(req *request.OnlineCheckinRequest, flightInfo *dto.FlightBasicInfo) ([]*dto.BoardingPassInfo, *response.Response) {
    var boardingPasses []*dto.BoardingPassInfo
    
    for _, checkinReq := range req.Checkins {
        log.Printf("Processing check-in for BookingDetailID: %d, SeatNumber: %s", 
            checkinReq.BookingDetailID, checkinReq.SeatNumber)
        
        // 1. Validate booking detail ownership
        passengerInfo, err := s.checkinRepo.ValidateBookingDetailOwnership(checkinReq.BookingDetailID, req.BookingID)
        if err != nil {
            log.Printf("Ownership validation failed for BookingDetailID %d: %v", checkinReq.BookingDetailID, err)
            return nil, &response.Response{
                Status:       false,
                ErrorCode:    error_code.InvalidRequest,
                ErrorMessage: fmt.Sprintf("Booking detail %d does not belong to this booking", checkinReq.BookingDetailID),
            }
        }
        
        // 2. Check if already checked in
        alreadyCheckedIn, err := s.checkinRepo.CheckAlreadyCheckedIn(checkinReq.BookingDetailID)
        if err != nil {
            log.Printf("Error checking if already checked in: %v", err)
            return nil, &response.Response{
                Status:       false,
                ErrorCode:    error_code.InternalError,
                ErrorMessage: "Failed to validate check-in status",
            }
        }
        
        if alreadyCheckedIn {
            return nil, &response.Response{
                Status:       false,
                ErrorCode:    error_code.InvalidRequest,
                ErrorMessage: fmt.Sprintf("Passenger %s has already checked in", passengerInfo.PassengerName),
            }
        }
        
        // 3. Check if seat is already occupied
        seatOccupied, err := s.checkinRepo.CheckSeatOccupied(flightInfo.FlightID, checkinReq.SeatNumber)
        if err != nil {
            log.Printf("Error checking seat occupancy: %v", err)
            return nil, &response.Response{
                Status:       false,
                ErrorCode:    error_code.InternalError,
                ErrorMessage: "Failed to check seat availability",
            }
        }
        
        if seatOccupied {
            return nil, &response.Response{
                Status:       false,
                ErrorCode:    error_code.InvalidRequest,
                ErrorMessage: fmt.Sprintf("Seat %s is already occupied", checkinReq.SeatNumber),
            }
        }
        
        // 4. Assign seat to passenger
        err = s.checkinRepo.AssignSeatToBookingDetail(checkinReq.BookingDetailID, checkinReq.SeatNumber, flightInfo.FlightID)
        if err != nil {
            log.Printf("Error assigning seat: %v", err)
            return nil, &response.Response{
                Status:       false,
                ErrorCode:    error_code.InternalError,
                ErrorMessage: "Failed to assign seat",
            }
        }
        
        // 5. Create check-in record
        boardingPassCode := s.generateBoardingPassCode()
        checkinData := &dto.CheckinDTO{
            BookingDetailId:  checkinReq.BookingDetailID,
            FlightID:         flightInfo.FlightID,
            Status:          "checked_in",
            CheckinTIme:     time.Now(),
            BoardingPassCode: boardingPassCode,
            CreatedAt:       time.Now(),
            UpdateAt:        time.Now(),
        }
        
        createdCheckin, err := s.checkinRepo.CreateCheckinRecord(checkinData)
        if err != nil {
            log.Printf("Error creating check-in record: %v", err)
            return nil, &response.Response{
                Status:       false,
                ErrorCode:    error_code.InternalError,
                ErrorMessage: "Failed to create check-in record",
            }
        }
        
        // 6. Create boarding pass
        boardingPass := &dto.BoardingPassInfo{
            BookingDetailID:  checkinReq.BookingDetailID,
            PassengerName:    passengerInfo.PassengerName,
            BoardingPassCode: createdCheckin.BoardingPassCode,
            SeatNumber:       checkinReq.SeatNumber,
            FlightClassName:  passengerInfo.FlightClassName,
            CheckinTime:      createdCheckin.CheckinTIme,
            Status:           createdCheckin.Status,
        }
        
        boardingPasses = append(boardingPasses, boardingPass)
        
        log.Printf("Successfully checked in passenger %s to seat %s", 
            passengerInfo.PassengerName, checkinReq.SeatNumber)
    }
    
    return boardingPasses, nil
}

// buildCheckinResponse - Xây dựng response cho check-in thành công
func (s *checkinService) buildCheckinResponse(bookingID int64, flightInfo *dto.FlightBasicInfo, boardingPasses []*dto.BoardingPassInfo) *response.Response {
    // Get PNR code
    _, err := s.bookingRepo.GetBookingByID(bookingID)
    if err != nil {
        log.Printf("Error getting booking info: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Failed to get booking information",
        }
    }
    var boardingPassValues []dto.BoardingPassInfo
    for _, bp := range boardingPasses {
        boardingPassValues = append(boardingPassValues, *bp)
    }
    // Calculate boarding time (usually 30 minutes before departure)
    boardingTime := flightInfo.DepartureTime.Add(-30 * time.Minute)
    
    checkinResponse := &dto.OnlineCheckinResponse{
        BookingID:       bookingID,
        
        FlightNumber:    flightInfo.FlightNumber,
        
        CheckinTime:     time.Now(),
        BoardingTime:    &boardingTime,
        CheckedInCount:  len(boardingPasses),
        BoardingPasses:  boardingPassValues,
    }
    
    log.Printf("Check-in completed successfully for booking %d, %d passengers checked in", 
        bookingID, len(boardingPasses))
    
    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: fmt.Sprintf("Successfully checked in %d passengers", len(boardingPasses)),
        Data:         checkinResponse,
    }
}


// generateBoardingPassCode - Tạo mã boarding pass
func (s *checkinService) generateBoardingPassCode() string {
    return fmt.Sprintf("BP%d%d", time.Now().Unix()%10000, time.Now().Nanosecond()%1000)
}