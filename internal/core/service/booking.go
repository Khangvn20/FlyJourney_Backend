package service

import (
    "fmt"
    "log"
    "time"
    "database/sql"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/dto"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/entity/error_code"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/model/request"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/model/response"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/port/repository"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/port/service"
)
type bookingService struct {
    bookingRepo  repository.BookingRepository
    redisService service.RedisService  
    pnrRepo      repository.PnrRepository 
}
func NewBookingService(bookingRepo repository.BookingRepository, redisService service.RedisService, pnrRepo repository.PnrRepository) service.BookingService {
    return &bookingService{
        bookingRepo:  bookingRepo,
        redisService: redisService,
        pnrRepo:      pnrRepo,
    }
}

//Rate limit booking
func (s *bookingService) RateLimitBooking(userID int64) *response.Response {
	key := fmt.Sprintf("booking_rate_limit:%d", userID)
    limit := 5
    ttl := 180
    count, err :=s.redisService.Incr(key)
    if err !=nil{
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: fmt.Sprintf("Error incrementing rate limit: %v", err),
        }
    }
    if count == 1 {
        if err :=s.redisService.Expire(key, time.Duration(ttl) * time.Second); err != nil {
            return &response.Response{
                Status:       false,
                ErrorCode:    error_code.InternalError,
                ErrorMessage: fmt.Sprintf("Error setting rate limit: %v", err),
            }
        }
    }

     if count > int64(limit) {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.RateLimitExceeded,
            ErrorMessage: fmt.Sprintf("Rate limit exceeded: %d requests per minute", limit),
        }
    }

    return nil
}
func (s *bookingService) CreateBooking(req *request.CreateBookingRequest) *response.Response {
    //Check rate limit
    rateLimitResponse := s.RateLimitBooking(req.UserID)
    if rateLimitResponse != nil {
        return rateLimitResponse
    }

	var lockedKeys []string
    lockValue := fmt.Sprintf("user:%d", req.UserID)

    passengerCountByClass := make(map[int64]int)
	for _, detail := range req.Details {
        passengerCountByClass[detail.FlightClassID]++
        if detail.ReturnFlightClassID != nil {
            passengerCountByClass[*detail.ReturnFlightClassID]++
        }
    }
 for flightClassID := range passengerCountByClass {
        lockKey := fmt.Sprintf("booking_lock:class:%d", flightClassID)
        locked, err := s.redisService.TryLock(lockKey, lockValue, 30*time.Second)

        if err != nil {

            for _, key := range lockedKeys {
                s.redisService.ReleaseLock(key, lockValue)
            }
            
            log.Printf("Failed to acquire lock for flight class %d: %v", flightClassID, err)
            return &response.Response{
                Status:       false,
                ErrorCode:    error_code.InternalError,
                ErrorMessage: "Hệ thống hiện không khả dụng",
            }
        }
        
        if !locked {
            // Giải phóng các khóa đã lấy được trước đó
            for _, key := range lockedKeys {
                s.redisService.ReleaseLock(key, lockValue)
            }
            
            return &response.Response{
                Status:       false,
                ErrorCode:    error_code.ResourceLocked,
                ErrorMessage: fmt.Sprintf("Hạng ghế %d đang được đặt bởi người khác, vui lòng thử lại sau", flightClassID),
            }
        }
        
        lockedKeys = append(lockedKeys, lockKey)
    }
    
    // 4. Đảm bảo giải phóng tất cả khóa khi hàm kết thúc
    defer func() {
        for _, key := range lockedKeys {
            err := s.redisService.ReleaseLock(key, lockValue)
            if err != nil {
                log.Printf("Error releasing lock %s: %v", key, err)
            }
        }
    }()
    
    // 5. Kiểm tra số ghế trống cho từng hạng ghế
    for flightClassID, passengerCount := range passengerCountByClass {
        available, seats, err := s.bookingRepo.CheckFlightClassAvailability(flightClassID)
        if err != nil {
            return &response.Response{
                Status:       false,
                ErrorCode:    error_code.InternalError,
                ErrorMessage: fmt.Sprintf("Lỗi kiểm tra chỗ trống: %v", err),
            }
        }
        
        if !available {
            return &response.Response{
                Status:       false,
                ErrorCode:    error_code.NoAvailableSeats,
                ErrorMessage: fmt.Sprintf("Hạng ghế %d đã hết chỗ", flightClassID),
            }
        }
        
        if seats < passengerCount {
            return &response.Response{
                Status:       false,
                ErrorCode:    error_code.InsufficientSeats,
                ErrorMessage: fmt.Sprintf("Hạng ghế %d chỉ còn %d chỗ trống (cần %d chỗ)", 
                                         flightClassID, seats, passengerCount),
            }
        }
    }
    
    // 6. Tạo booking object
    booking := &dto.Booking{
        UserID:         req.UserID,
        FlightID:       req.FlightID,
        ReturnFlightID: req.ReturnFlightID,
        ContactEmail:   req.ContactEmail,
        ContactPhone:   req.ContactPhone,
        ContactAddress: req.ContactAddress,
        Note:           req.Note,
        TotalPrice:     req.TotalPrice,
        Status:         "pending_payment",
        BookingDate:    time.Now(),
        CheckInStatus:  "not_checked_in",
        Details:        make([]*dto.BookingDetail, len(req.Details)),
    }
    
    // 7. Xử lý booking details - chuyển đổi thông tin hành khách
    for i, detail := range req.Details {
        // Chuyển đổi chuỗi ngày thành time.Time
        dob, err := time.Parse("02/01/2006", detail.DateOfBirth)
        if err != nil {
            return &response.Response{
                Status:       false,
                ErrorCode:    error_code.InvalidRequest,
                ErrorMessage: fmt.Sprintf("Ngày sinh không hợp lệ: %s", detail.DateOfBirth),
            }
        }
        
        expiry, err := time.Parse("02/01/2006", detail.ExpiryDate)
        if err != nil {
            return &response.Response{
                Status:       false,
                ErrorCode:    error_code.InvalidRequest,
                ErrorMessage: fmt.Sprintf("Ngày hết hạn không hợp lệ: %s", detail.ExpiryDate),
            }
        }
        
        booking.Details[i] = &dto.BookingDetail{
            PassengerAge:    detail.PassengerAge,
            PassengerGender: detail.PassengerGender,
            FlightClassID:   detail.FlightClassID,
            ReturnFlightClassID: detail.ReturnFlightClassID,
            Price:           detail.Price,
            LastName:        detail.LastName,
            FirstName:       detail.FirstName,
            DateOfBirth:     dob,
            IDType:          detail.IDType,
            IDNumber:        detail.IDNumber,
            ExpiryDate:      expiry,
            IssuingCountry:  detail.IssuingCountry,
            Nationality:     detail.Nationality,
        }
    }
    
    // 8. Xử lý ancillaries - dịch vụ bổ sung
    if req.Ancillaries != nil && len(req.Ancillaries) > 0 {
        booking.Ancillaries = make([]*dto.Ancillary, len(req.Ancillaries))
        for i, ancillary := range req.Ancillaries {
            booking.Ancillaries[i] = &dto.Ancillary{
                Type:        ancillary.Type,
                Description: ancillary.Description,
                Quantity:    ancillary.Quantity,
                Price:       ancillary.Price,
                CreatedAt:   time.Now(),
            }
        }
    }
  
    createdBooking, err := s.bookingRepo.CreateBooking(booking)
    if err != nil {
        log.Printf("Error creating booking: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: fmt.Sprintf("Lỗi tạo đặt chỗ: %v", err),
        }
    }
    pnr := &dto.PNR{
        BookingID:      fmt.Sprintf("%d", createdBooking.BookingID),
        FlightID:       fmt.Sprintf("%d", createdBooking.FlightID),
        ReturnFlightID: sql.NullString{},
        Status:         "active",
    }

    // Nếu có chuyến về, thêm vào PNR
    if createdBooking.ReturnFlightID != nil {
        pnr.ReturnFlightID = sql.NullString{
            String: fmt.Sprintf("%d", *createdBooking.ReturnFlightID),
            Valid:  true,
        }
    }

    // Tạo PNR record
    createdPNR, err := s.pnrRepo.CreatePnr(pnr)
    if err != nil {
        log.Printf("Error creating PNR for booking %d: %v", createdBooking.BookingID, err)
    } else {
        log.Printf("PNR %s created successfully for booking %d", createdPNR.PNRCode, createdBooking.BookingID)
    }
    // 10. Đặt timeout cho booking (2 giờ)
    timeoutKey := fmt.Sprintf("booking:timeout:%d", createdBooking.BookingID)
    if err := s.redisService.Set(timeoutKey, "pending_payment", 2*time.Hour); err != nil {
        log.Printf("Warning: Could not set booking timeout: %v", err)
    }
    
    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Đặt chỗ thành công",
        Data:         createdBooking,
    }
}
//Worker 2h expire bookings
func (s *bookingService) CancelExpiredBookings() *response.Response {
    expiredBookingIDs, err := s.bookingRepo.GetExpiredBookingIDs()
    if err != nil {
        log.Printf("Error fetching expired booking IDs: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: fmt.Sprintf("Lỗi lấy danh sách các đặt chỗ quá hạn: %v", err),
        }
    }

    if len(expiredBookingIDs) == 0 {
        log.Println("No expired bookings found to cancel")
        return &response.Response{
            Status:       true,
            ErrorCode:    error_code.Success,
            ErrorMessage: "Không có đặt chỗ quá hạn để hủy",
        }
    }
    canceledBookingIDs, err := s.bookingRepo.CancelBookings(expiredBookingIDs)
    if err != nil {
        log.Printf("Error canceling expired bookings: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: fmt.Sprintf("Lỗi hủy các đặt chỗ quá hạn: %v", err),
        }
    }

    log.Printf("Canceled %d expired bookings: %v", len(canceledBookingIDs), canceledBookingIDs)
    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: fmt.Sprintf("Hủy thành công %d đặt chỗ quá hạn", len(canceledBookingIDs)),
        Data:         canceledBookingIDs,
    }
}