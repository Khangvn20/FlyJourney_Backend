package service
import (
    "bytes"
    "fmt"
    "html/template"
    "log"
    "time"  
    "github.com/dustin/go-humanize"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/dto"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/entity/error_code"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/model/response"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/port/repository"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/port/service"
)
type bookingEmailService struct {
	bookingRepo   repository.BookingRepository
    flightRepo    repository.FlightRepository
    pnrRepo       repository.PnrRepository
    userRepo      repository.UserRepository
    paymentRepo   repository.PaymentRepository
    emailService  service.EmailOTPService 
}
func NewBookingEmailService(
    bookingRepo repository.BookingRepository,
    flightRepo repository.FlightRepository,
    pnrRepo repository.PnrRepository,
    userRepo repository.UserRepository,
    paymentRepo repository.PaymentRepository,
    emailService service.EmailOTPService,
) service.BookingEmailService {
    return &bookingEmailService{
        bookingRepo:  bookingRepo,
        flightRepo:   flightRepo,
        pnrRepo:      pnrRepo,
        userRepo:     userRepo,
        paymentRepo:  paymentRepo,
        emailService: emailService,
    }
}
func (s *bookingEmailService) SendBookingConfirmationEmail(bookingID int64) *response.Response {

	booking, err := s.bookingRepo.GetBookingByID(bookingID)
    if err != nil {
        log.Printf("Error fetching booking: %v", err)
        return &response.Response{
            Status: false,
            ErrorCode: error_code.InternalError,
            ErrorMessage: "fail to fetch booking ID",
        }
    }
	if booking == nil {
		return &response.Response{
			Status: false,
			ErrorCode: error_code.NOTFOUND,
			ErrorMessage: "Booking not found",
		}
	}
	//get user by id
	user, err := s.userRepo.GetUserByID(int(booking.UserID))
    if err != nil {
        log.Printf("Error fetching user: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "fail to fetch user ID",
        }
    }

	pnr, err := s.pnrRepo.GetPnrByBookingID(bookingID)
	if err != nil {
		log.Printf("Error fetching PNR: %v", err)
		return &response.Response{
			Status: false,
			ErrorCode: error_code.InternalErrMsg,
			ErrorMessage: "fail to fetch PNR",
		}
	}
	outboundFlight, _, err := s.flightRepo.GetByID(int(booking.FlightID))
    if err != nil {
    log.Printf("Error fetching outbound flight: %v", err)
    return &response.Response{
        Status: false,
        ErrorCode: error_code.InternalErrMsg,
        ErrorMessage: "fail to fetch outbound flight",
    }
}
	payment, err := s.paymentRepo.GetPaymentByBookingID(bookingID)
	if err != nil {
		log.Printf("Error fetching payment: %v", err)
		return &response.Response{
			Status: false,
			ErrorCode: error_code.InternalErrMsg,
			ErrorMessage: "fail to fetch payment",
		}
	}
	if payment.Status != "success" {
        log.Printf("Payment not completed for booking ID: %d", bookingID)
        return &response.Response{
            Status: false,
            ErrorCode: error_code.NOTFOUND,
            ErrorMessage: "Payment not completed",
        }
    }

    emailData :=&dto.BookingEmailData{
		
		 PNRCode:      pnr.PNRCode,
        BookingID:    booking.BookingID,
        UserFullName: user.Name,
        ContactEmail: booking.ContactEmail,
        ContactPhone: booking.ContactPhone,
        TotalPrice:  humanize.Comma(int64(booking.TotalPrice)),
        PaymentDate:  time.Now(),
        OutboundFlight: &dto.BookingEmailFlight{
            FlightNumber:     outboundFlight.FlightNumber,
            AirlineName:      outboundFlight.AirlineName,
            DepartureAirport: fmt.Sprintf("%s (%s)", outboundFlight.DepartureAirport, outboundFlight.DepartureAirportCode),
            ArrivalAirport:   fmt.Sprintf("%s (%s)", outboundFlight.ArrivalAirport, outboundFlight.ArrivalAiportCode),
            DepartureTime:    outboundFlight.DepartureTime,
            ArrivalTime:      outboundFlight.ArrivalTime,
        },
        Passengers: make([]*dto.BookingEmailPassenger, len(booking.Details)),
    }
	if payment != nil {
		emailData.PaymentMethod = payment.PaymentMethod
		emailData.PaymentDate = *payment.PaidAt
	}
	if booking.ReturnFlightID != nil && *booking.ReturnFlightID > 0 {
        returnFlight, _, err := s.flightRepo.GetByID(int(*booking.ReturnFlightID))
        if err == nil && returnFlight != nil {
            emailData.InboundFlight = &dto.BookingEmailFlight{
                FlightNumber:     returnFlight.FlightNumber,
                AirlineName:      returnFlight.AirlineName,
                DepartureAirport: fmt.Sprintf("%s (%s)", returnFlight.DepartureAirport, returnFlight.DepartureAirportCode),
                ArrivalAirport:   fmt.Sprintf("%s (%s)", returnFlight.ArrivalAirport, returnFlight.ArrivalAiportCode),
                DepartureTime:    returnFlight.DepartureTime,
                ArrivalTime:      returnFlight.ArrivalTime,
            }
        }
    }
	for i, detail := range booking.Details {
		passengerType :="Adult"
		if detail.PassengerAge <12 && detail.PassengerAge >=2 {
			passengerType = "Child"
		}else if detail.PassengerAge < 2 {
			passengerType = "Infant"
		}
		emailData.Passengers[i] = &dto.BookingEmailPassenger{
            FullName:      fmt.Sprintf("%s %s", detail.FirstName, detail.LastName),
            Type: passengerType,
            FlightClass: detail.FlightClassName,
      // optimize seat number
            SeatNumber: "",
        }
    }
    log.Printf("Booking Details: %+v", booking.Details)
    htmlContent, err := s.renderEmailTemplate(emailData)
    if err != nil {
        log.Printf("Error rendering email template: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Không thể tạo nội dung email",
        }
    }
    subject := fmt.Sprintf("Xác nhận đặt vé - Thông tin chi tiết vé của bạn")
	 err = s.emailService.(*emailOTPService).SendHTMLMail(booking.ContactEmail, subject, htmlContent)
    if err != nil {
        log.Printf("Error sending confirmation email: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Không thể gửi email xác nhận",
        }
    }

    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Đã gửi email xác nhận đặt vé thành công",
    }
}

func (s *bookingEmailService) renderEmailTemplate(data *dto.BookingEmailData) (string, error) {
 
const tmplStr = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Xác nhận đặt vé máy bay</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 800px; margin: 0 auto; }
        .header { background-color: #0066cc; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; }
        .booking-info { margin-bottom: 20px; border: 1px solid #ddd; padding: 15px; background-color: #f9f9f9; }
        .flight-info { background-color: #f5f5f5; padding: 15px; margin-bottom: 20px; border-left: 4px solid #0066cc; }
        .passengers { margin-bottom: 20px; }
        table { width: 100%; border-collapse: collapse; }
        table, th, td { border: 1px solid #ddd; }
        th, td { padding: 10px; text-align: left; }
        th { background-color: #f2f2f2; }
        .footer { background-color: #f2f2f2; padding: 15px; text-align: center; font-size: 14px; }
        .booking-code { font-size: 24px; font-weight: bold; color: #0066cc; }
        .label { font-weight: bold; }
        .total-price { font-size: 20px; color: #d9534f; font-weight: bold; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Xác nhận đặt vé chuyến bay</h1>
    </div>
    
    <div class="content">
        <p>Kính chào <strong>{{.UserFullName}}</strong>,</p>
        <p>Cảm ơn bạn đã đặt vé tại FlyJourney. Đơn đặt vé của bạn đã được thanh toán thành công. Dưới đây là thông tin chi tiết về chuyến bay của bạn:</p>
        
        <div class="booking-info">
            <p><span class="label">Mã đặt chỗ (PNR):</span> <span class="booking-code">{{.PNRCode}}</span></p>
            <p><span class="label">Mã đơn hàng:</span> #{{.BookingID}}</p>
            <p><span class="label">Ngày đặt:</span> {{.PaymentDate.Format "02/01/2006 15:04"}}</p>
            <p><span class="label">Phương thức thanh toán:</span> {{.PaymentMethod}}</p>
           <p><span class="label">Tổng thanh toán:</span> <span class="total-price">{{.TotalPrice}} VND</span></p>
        </div>
        
        <h2>Thông tin chuyến bay</h2>

        {{if .OutboundFlight}}
        <div class="flight-info">
            <h3>Chuyến đi</h3>
            <p><span class="label">Hãng bay:</span> {{.OutboundFlight.AirlineName}}</p>
            <p><span class="label">Số hiệu chuyến bay:</span> {{.OutboundFlight.FlightNumber}}</p>
            <p><span class="label">Khởi hành:</span> {{.OutboundFlight.DepartureAirport}}</p>
            <p><span class="label">Thời gian:</span> {{.OutboundFlight.DepartureTime.Format "02/01/2006 15:04"}}</p>
            <p><span class="label">Đến:</span> {{.OutboundFlight.ArrivalAirport}}</p>
            <p><span class="label">Thời gian:</span> {{.OutboundFlight.ArrivalTime.Format "02/01/2006 15:04"}}</p>
            <p><span class="label">Hạng vé:</span> {{.OutboundFlight.FlightClass}}</p>
        </div>
        {{end}}

        {{if .InboundFlight}}
        <div class="flight-info">
            <h3>Chuyến về</h3>
            <p><span class="label">Hãng bay:</span> {{.InboundFlight.AirlineName}}</p>
            <p><span class="label">Số hiệu chuyến bay:</span> {{.InboundFlight.FlightNumber}}</p>
            <p><span class="label">Khởi hành:</span> {{.InboundFlight.DepartureAirport}}</p>
            <p><span class="label">Thời gian:</span> {{.InboundFlight.DepartureTime.Format "02/01/2006 15:04"}}</p>
            <p><span class="label">Đến:</span> {{.InboundFlight.ArrivalAirport}}</p>
            <p><span class="label">Thời gian:</span> {{.InboundFlight.ArrivalTime.Format "02/01/2006 15:04"}}</p>
            <p><span class="label">Hạng vé:</span> {{.InboundFlight.FlightClass}}</p>
        </div>
        {{end}}
        
        <div class="passengers">
            <h2>Thông tin hành khách</h2>
            <table>
                <tr>
                    <th>Họ tên</th>
                    <th>Loại</th>
                    <th>Hạng vé</th>
                    <th>Chỗ ngồi</th>
                </tr>
                {{range .Passengers}}
                <tr>
                    <td>{{.FullName}}</td>
                    <td>{{.Type}}</td>
                    <td>{{.FlightClass}}</td>
                    <td>{{if .SeatNumber}}{{.SeatNumber}}{{else}}Chưa chọn{{end}}</td>
                </tr>
                {{end}}
            </table>
        </div>
        
        <p>Vui lòng kiểm tra kỹ thông tin và liên hệ với chúng tôi nếu có bất kỳ câu hỏi nào.</p>
        <p>Để check-in trực tuyến, vui lòng truy cập website hoặc ứng dụng FlyJourney trước giờ khởi hành 24 giờ.</p>
        <p>Bạn nên có mặt tại sân bay ít nhất 2 giờ trước giờ khởi hành đối với chuyến bay nội địa và 3 giờ đối với chuyến bay quốc tế.</p>
    </div>
    
    <div class="footer">
        <p>Đây là email tự động. Vui lòng không trả lời email này.</p>
        <p>&copy; 2025 FlyJourney - Hệ thống đặt vé máy bay trực tuyến</p>
        <p>Nếu bạn cần hỗ trợ, vui lòng liên hệ: <a href="mailto:support@flyjourney.com">support@flyjourney.com</a></p>
    </div>
</body>
</html>`
    // Tạo template
    tmpl, err := template.New("bookingConfirmation").Parse(tmplStr)
    if err != nil {
        return "", err
    }

    // Render template với dữ liệu
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", err
    }

    return buf.String(), nil
}

func (s *bookingEmailService) SendFlightCancelEmail(bookingID int64, reason string) *response.Response {
    booking, err := s.bookingRepo.GetBookingByID(bookingID)
    if err != nil {
        log.Printf("error fetch booking: %v", err)
        return &response.Response{
            Status: false,
            ErrorCode: error_code.InternalError,
            ErrorMessage: "error fetching booking",
        }
    }
   user, err :=s.userRepo.GetUserByID(int(booking.UserID))
    if err != nil {
        log.Printf("error fetch user: %v", err)
        return &response.Response{
            Status: false,
            ErrorCode: error_code.InternalError,
            ErrorMessage: "error fetching user",
        }
    }
    pnr, err := s.pnrRepo.GetPnrByBookingID(bookingID)
    if err != nil {
        log.Printf("error fetch pnr: %v", err)
        return &response.Response{
            Status: false,
            ErrorCode: error_code.InternalError,
            ErrorMessage: "error fetching pnr",
        }
    }
    outboundFlight, _, err := s.flightRepo.GetByID(int(booking.FlightID))
    if err != nil {
        log.Printf("Error fetching outbound flight: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Không thể lấy thông tin chuyến bay",
        }
    }
    emailData := &dto.CancellationEmailData{
        PNRCode:            pnr.PNRCode,
        BookingID:          booking.BookingID,
        UserFullName:       user.Name,
        ContactEmail:       booking.ContactEmail,
        ContactPhone:       booking.ContactPhone,
        CancellationReason: reason,
        RefundAmount:       booking.TotalPrice,
        CancellationDate:   time.Now(),
        OutboundFlight: &dto.BookingEmailFlight{
            FlightNumber:     outboundFlight.FlightNumber,
            AirlineName:      outboundFlight.AirlineName,
            DepartureAirport: fmt.Sprintf("%s (%s)", outboundFlight.DepartureAirport, outboundFlight.DepartureAirportCode),
            ArrivalAirport:   fmt.Sprintf("%s (%s)", outboundFlight.ArrivalAirport, outboundFlight.ArrivalAiportCode),
            DepartureTime:    outboundFlight.DepartureTime,
            ArrivalTime:      outboundFlight.ArrivalTime,
        },
    }
    if booking.ReturnFlightID != nil && *booking.ReturnFlightID > 0 {
        returnFlight, _, err := s.flightRepo.GetByID(int(*booking.ReturnFlightID))
        if err == nil && returnFlight != nil {
            emailData.InboundFlight = &dto.BookingEmailFlight{
                FlightNumber:     returnFlight.FlightNumber,
                AirlineName:      returnFlight.AirlineName,
                DepartureAirport: fmt.Sprintf("%s (%s)", returnFlight.DepartureAirport, returnFlight.DepartureAirportCode),
                ArrivalAirport:   fmt.Sprintf("%s (%s)", returnFlight.ArrivalAirport, returnFlight.ArrivalAiportCode),
                DepartureTime:    returnFlight.DepartureTime,
                ArrivalTime:      returnFlight.ArrivalTime,
            }
        }
    }
        htmlContent, err := s.renderCancellationEmailTemplate(emailData)
    if err != nil {
        log.Printf("Error rendering cancellation email template: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Không thể tạo nội dung email",
        }
    }
    
 
    subject := "Thông báo hủy chuyến bay - Cập nhật thông tin đặt vé của bạn"
    err = s.emailService.(*emailOTPService).SendHTMLMail(booking.ContactEmail, subject, htmlContent)
    if err != nil {
        log.Printf("Error sending cancellation email: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Không thể gửi email thông báo",
        }
    }
    
    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Đã gửi email thông báo hủy chuyến bay thành công",
    }
}


func (s *bookingEmailService) renderCancellationEmailTemplate(data *dto.CancellationEmailData) (string, error) {
    tmplStr := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Thông báo hủy chuyến bay</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 800px; margin: 0 auto; }
        .header { background-color: #d9534f; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; }
        .booking-info { margin-bottom: 20px; border: 1px solid #ddd; padding: 15px; background-color: #f9f9f9; }
        .flight-info { background-color: #f5f5f5; padding: 15px; margin-bottom: 20px; border-left: 4px solid #d9534f; }
        .refund-info { margin-top: 20px; padding: 15px; background-color: #dff0d8; border: 1px solid #d6e9c6; }
        .footer { background-color: #f2f2f2; padding: 15px; text-align: center; font-size: 14px; }
        .label { font-weight: bold; }
        .reason { color: #d9534f; font-weight: bold; }
        .alert { color: #d9534f; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Thông báo hủy chuyến bay</h1>
    </div>
    
    <div class="content">
        <p>Kính chào <strong>{{.UserFullName}}</strong>,</p>
        <p>Chúng tôi rất tiếc phải thông báo rằng chuyến bay của bạn đã bị hủy. Dưới đây là thông tin chi tiết:</p>
        
        <div class="booking-info">
            <p><span class="label">Mã đặt chỗ (PNR):</span> {{.PNRCode}}</p>
            <p><span class="label">Mã đơn hàng:</span> #{{.BookingID}}</p>
            <p><span class="label">Lý do hủy chuyến:</span> <span class="reason">{{.CancellationReason}}</span></p>
        </div>
        
        <h2>Thông tin chuyến bay bị hủy</h2>

        <div class="flight-info">
            <h3>Chuyến đi</h3>
            <p><span class="label">Hãng bay:</span> {{.OutboundFlight.AirlineName}}</p>
            <p><span class="label">Số hiệu chuyến bay:</span> {{.OutboundFlight.FlightNumber}}</p>
            <p><span class="label">Khởi hành:</span> {{.OutboundFlight.DepartureAirport}}</p>
            <p><span class="label">Thời gian:</span> {{.OutboundFlight.DepartureTime.Format "02/01/2006 15:04"}}</p>
            <p><span class="label">Đến:</span> {{.OutboundFlight.ArrivalAirport}}</p>
            <p><span class="label">Thời gian:</span> {{.OutboundFlight.ArrivalTime.Format "02/01/2006 15:04"}}</p>
        </div>
        
        {{if .InboundFlight}}
        <div class="flight-info">
            <h3>Chuyến về</h3>
            <p><span class="label">Hãng bay:</span> {{.InboundFlight.AirlineName}}</p>
            <p><span class="label">Số hiệu chuyến bay:</span> {{.InboundFlight.FlightNumber}}</p>
            <p><span class="label">Khởi hành:</span> {{.InboundFlight.DepartureAirport}}</p>
            <p><span class="label">Thời gian:</span> {{.InboundFlight.DepartureTime.Format "02/01/2006 15:04"}}</p>
            <p><span class="label">Đến:</span> {{.InboundFlight.ArrivalAirport}}</p>
            <p><span class="label">Thời gian:</span> {{.InboundFlight.ArrivalTime.Format "02/01/2006 15:04"}}</p>
        </div>
        {{end}}
        
        <div class="refund-info">
            <h3>Thông tin hoàn tiền</h3>
            <p>Theo chính sách của chúng tôi, bạn sẽ được hoàn lại số tiền: <strong>{{.RefundAmount}} VND</strong></p>
            <p>Tiền hoàn lại sẽ được chuyển về phương thức thanh toán ban đầu của bạn trong vòng 7-14 ngày làm việc.</p>
        </div>
        
        <p class="alert">Chúng tôi thành thật xin lỗi vì sự bất tiện này và đang cố gắng hỗ trợ bạn tốt nhất có thể.</p>
        
        <p>Nếu bạn muốn đặt lại chuyến bay hoặc cần hỗ trợ thêm, vui lòng liên hệ với đội ngũ hỗ trợ khách hàng của chúng tôi qua:</p>
        <ul>
            <li>Hotline: 1900 xxxx</li>
            <li>Email: support@flyjourney.com</li>
            <li>Hoặc chat trực tiếp trên website/ứng dụng FlyJourney</li>
        </ul>
    </div>
    
    <div class="footer">
        <p>Đây là email tự động. Vui lòng không trả lời email này.</p>
        <p>&copy; 2025 FlyJourney - Hệ thống đặt vé máy bay trực tuyến</p>
    </div>
</body>
</html>`

    // Tạo template
    tmpl, err := template.New("flightCancellation").Parse(tmplStr)
    if err != nil {
        return "", err
    }

    // Render template với dữ liệu
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", err
    }

    return buf.String(), nil
} 


func (s *bookingEmailService) SendFlightDelayEmail(bookingID int64, newDepartureTime time.Time, reason string) *response.Response {
    // Lấy thông tin booking
    booking, err := s.bookingRepo.GetBookingByID(bookingID)
    if err != nil {
        log.Printf("error fetch booking: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "error fetching booking",
        }
    }

    // Lấy thông tin user
    user, err := s.userRepo.GetUserByID(int(booking.UserID))
    if err != nil {
        log.Printf("error fetch user: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "error fetching user",
        }
    }

    // Lấy thông tin PNR
    pnr, err := s.pnrRepo.GetPnrByBookingID(bookingID)
    if err != nil {
        log.Printf("error fetch pnr: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "error fetching pnr",
        }
    }

    // Lấy thông tin chuyến bay
    outboundFlight, _, err := s.flightRepo.GetByID(int(booking.FlightID))
    if err != nil {
        log.Printf("Error fetching outbound flight: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Không thể lấy thông tin chuyến bay",
        }
    }


    originalTime := outboundFlight.DepartureTime
    delayDuration := newDepartureTime.Sub(originalTime)

    var delayDurationStr string
    if delayDuration.Hours() >= 1 {
        hours := int(delayDuration.Hours())
        minutes := int(delayDuration.Minutes()) % 60
        if minutes > 0 {
            delayDurationStr = fmt.Sprintf("%d giờ %d phút", hours, minutes)
        } else {
            delayDurationStr = fmt.Sprintf("%d giờ", hours)
        }
    } else {
        delayDurationStr = fmt.Sprintf("%d phút", int(delayDuration.Minutes()))
    }

    emailData := &dto.DelayEmailData{
        PNRCode:          pnr.PNRCode,
        BookingID:        booking.BookingID,
        UserFullName:     user.Name,
        ContactEmail:     booking.ContactEmail,
        ContactPhone:     booking.ContactPhone,
        DelayReason:      reason,
        OriginalTime:     originalTime,
        NewDepartureTime: newDepartureTime,
        DelayDuration:    delayDurationStr,
        NotificationTime: time.Now(),
        OutboundFlight: &dto.BookingEmailFlight{
            FlightNumber:     outboundFlight.FlightNumber,
            AirlineName:      outboundFlight.AirlineName,
            DepartureAirport: fmt.Sprintf("%s (%s)", outboundFlight.DepartureAirport, outboundFlight.DepartureAirportCode),
            ArrivalAirport:   fmt.Sprintf("%s (%s)", outboundFlight.ArrivalAirport, outboundFlight.ArrivalAiportCode),
            DepartureTime:    newDepartureTime,
            ArrivalTime:      outboundFlight.ArrivalTime.Add(delayDuration), 
        },
    }

    if booking.ReturnFlightID != nil && *booking.ReturnFlightID > 0 {
        returnFlight, _, err := s.flightRepo.GetByID(int(*booking.ReturnFlightID))
        if err == nil && returnFlight != nil {
            emailData.InboundFlight = &dto.BookingEmailFlight{
                FlightNumber:     returnFlight.FlightNumber,
                AirlineName:      returnFlight.AirlineName,
                DepartureAirport: fmt.Sprintf("%s (%s)", returnFlight.DepartureAirport, returnFlight.DepartureAirportCode),
                ArrivalAirport:   fmt.Sprintf("%s (%s)", returnFlight.ArrivalAirport, returnFlight.ArrivalAiportCode),
                DepartureTime:    returnFlight.DepartureTime,
                ArrivalTime:      returnFlight.ArrivalTime,
            }
        }
    }

    // Render email template
    htmlContent, err := s.renderDelayEmailTemplate(emailData)
    if err != nil {
        log.Printf("Error rendering delay email template: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Không thể tạo nội dung email",
        }
    }


    subject := "Thông báo trì hoãn chuyến bay - Cập nhật thông tin đặt vé của bạn"
    err = s.emailService.(*emailOTPService).SendHTMLMail(booking.ContactEmail, subject, htmlContent)
    if err != nil {
        log.Printf("Error sending delay email: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Không thể gửi email thông báo",
        }
    }

    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Đã gửi email thông báo delay chuyến bay thành công",
    }
}
//render email Delay
func (s *bookingEmailService) renderDelayEmailTemplate(data *dto.DelayEmailData) (string, error) {
    tmplStr := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Thông báo trì hoãn chuyến bay</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 800px; margin: 0 auto; }
        .header { background-color: #f0ad4e; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; }
        .booking-info { margin-bottom: 20px; border: 1px solid #ddd; padding: 15px; background-color: #f9f9f9; }
        .flight-info { background-color: #f5f5f5; padding: 15px; margin-bottom: 20px; border-left: 4px solid #f0ad4e; }
        .delay-info { margin-top: 20px; padding: 15px; background-color: #fcf8e3; border: 1px solid #faebcc; }
        .footer { background-color: #f2f2f2; padding: 15px; text-align: center; font-size: 14px; }
        .label { font-weight: bold; }
        .delay-time { color: #f0ad4e; font-weight: bold; }
        .original-time { text-decoration: line-through; color: #777; }
        .new-time { font-weight: bold; color: #d9534f; }
        .alert { color: #8a6d3b; }
        table { width: 100%; border-collapse: collapse; margin-bottom: 20px; }
        th, td { padding: 10px; text-align: left; border-bottom: 1px solid #ddd; }
        th { background-color: #f2f2f2; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Thông báo trì hoãn chuyến bay</h1>
    </div>
    
    <div class="content">
        <p>Kính chào <strong>{{.UserFullName}}</strong>,</p>
        <p>Chúng tôi xin thông báo chuyến bay của bạn đã bị trì hoãn. Dưới đây là thông tin chi tiết:</p>
        
        <div class="booking-info">
            <p><span class="label">Mã đặt chỗ (PNR):</span> {{.PNRCode}}</p>
            <p><span class="label">Mã đơn hàng:</span> #{{.BookingID}}</p>
        </div>
        
        <div class="delay-info">
            <h3>Thông tin trì hoãn:</h3>
            <p><span class="label">Thời gian trì hoãn:</span> <span class="delay-time">{{.DelayDuration}}</span></p>
            <p><span class="label">Lý do:</span> {{.DelayReason}}</p>
            <p><span class="label">Thời gian khởi hành ban đầu:</span> <span class="original-time">{{.OriginalTime.Format "02/01/2006 15:04"}}</span></p>
            <p><span class="label">Thời gian khởi hành mới:</span> <span class="new-time">{{.NewDepartureTime.Format "02/01/2006 15:04"}}</span></p>
        </div>
        
        <h2>Chi tiết chuyến bay</h2>
        <div class="flight-info">
            <h3>Chuyến đi (Đã cập nhật)</h3>
            <table>
                <tr>
                    <th>Thông tin</th>
                    <th>Chi tiết</th>
                </tr>
                <tr>
                    <td>Hãng bay</td>
                    <td>{{.OutboundFlight.AirlineName}}</td>
                </tr>
                <tr>
                    <td>Số hiệu chuyến bay</td>
                    <td>{{.OutboundFlight.FlightNumber}}</td>
                </tr>
                <tr>
                    <td>Khởi hành từ</td>
                    <td>{{.OutboundFlight.DepartureAirport}}</td>
                </tr>
                <tr>
                    <td>Thời gian khởi hành mới</td>
                    <td class="new-time">{{.OutboundFlight.DepartureTime.Format "02/01/2006 15:04"}}</td>
                </tr>
                <tr>
                    <td>Đến</td>
                    <td>{{.OutboundFlight.ArrivalAirport}}</td>
                </tr>
                <tr>
                    <td>Thời gian đến (dự kiến)</td>
                    <td>{{.OutboundFlight.ArrivalTime.Format "02/01/2006 15:04"}}</td>
                </tr>
            </table>
        </div>
        
        {{if .InboundFlight}}
        <div class="flight-info">
            <h3>Chuyến về</h3>
            <table>
                <tr>
                    <th>Thông tin</th>
                    <th>Chi tiết</th>
                </tr>
                <tr>
                    <td>Hãng bay</td>
                    <td>{{.InboundFlight.AirlineName}}</td>
                </tr>
                <tr>
                    <td>Số hiệu chuyến bay</td>
                    <td>{{.InboundFlight.FlightNumber}}</td>
                </tr>
                <tr>
                    <td>Khởi hành từ</td>
                    <td>{{.InboundFlight.DepartureAirport}}</td>
                </tr>
                <tr>
                    <td>Thời gian khởi hành</td>
                    <td>{{.InboundFlight.DepartureTime.Format "02/01/2006 15:04"}}</td>
                </tr>
                <tr>
                    <td>Đến</td>
                    <td>{{.InboundFlight.ArrivalAirport}}</td>
                </tr>
                <tr>
                    <td>Thời gian đến</td>
                    <td>{{.InboundFlight.ArrivalTime.Format "02/01/2006 15:04"}}</td>
                </tr>
            </table>
        </div>
        {{end}}
        
        <div class="alert">
            <p>Chúng tôi rất tiếc về sự bất tiện này và đánh giá cao sự thông cảm của bạn.</p>
            <p>Vui lòng kiểm tra email và tin nhắn của bạn thường xuyên để nhận các cập nhật mới nhất về chuyến bay.</p>
        </div>
        
        <p>Nếu bạn có bất kỳ thắc mắc nào hoặc cần thay đổi kế hoạch, vui lòng liên hệ với đội ngũ hỗ trợ khách hàng của chúng tôi qua:</p>
        <ul>
            <li>Hotline: 1900 xxxx</li>
            <li>Email: support@flyjourney.com</li>
            <li>Hoặc chat trực tiếp trên website/ứng dụng FlyJourney</li>
        </ul>
    </div>
    
    <div class="footer">
        <p>Đây là email tự động. Vui lòng không trả lời email này.</p>
        <p>&copy; 2025 FlyJourney - Hệ thống đặt vé máy bay trực tuyến</p>
    </div>
</body>
</html>`



    tmpl, err := template.New("flightDelay").Parse(tmplStr)
    if err != nil {
        return "", err
    }

    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", err
    }

    return buf.String(), nil
}