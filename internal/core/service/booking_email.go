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