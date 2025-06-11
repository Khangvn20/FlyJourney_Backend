package service

import (
    "fmt"
    "math/rand"
    "time"
	"os"
	"net/smtp"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/model/response"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/entity/error_code"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/port/service"
)

type emailOTPService struct {
    otpRepo map[string]OTP 
}

type OTP struct {
    Email     string
    Code      string
    ExpiresAt time.Time
}

func NewEmailOTPService() service.EmailOTPService {
    return &emailOTPService{
        otpRepo: make(map[string]OTP),
    }
}
func sendMail(to, subject, body string) error {
      from := os.Getenv("email_user")
    password := os.Getenv("password_user")
    smtpHost := "smtp.gmail.com"
    smtpPort := "587"

    msg := "From: " + from + "\n" +
        "To: " + to + "\n" +
        "Subject: " + subject + "\n\n" +
        body

    auth := smtp.PlainAuth("", from, password, smtpHost)
    return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
}
func (s *emailOTPService) generateOTP() string {
    rand.Seed(time.Now().UnixNano())
    return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (s *emailOTPService) SendOTPEmail(email string) *response.Response {
    otpCode := s.generateOTP()
    expire := time.Now().Add(5 * time.Minute)
    s.otpRepo[email] = OTP{
        Email:     email,
        Code:      otpCode,
        ExpiresAt: expire,
    }
 
 subject := "Mã xác thực OTP của bạn"
    body := fmt.Sprintf("Mã OTP của bạn là: %s", otpCode)
    err := sendMail(email, subject, body)
    if err != nil {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Không thể gửi email OTP",
        }
    }
    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "OTP đã được gửi tới email",
    }
}

func (s *emailOTPService) VerifyEmail(email, otp string) *response.Response {
    stored, ok := s.otpRepo[email]
    if !ok || stored.Code != otp || time.Now().After(stored.ExpiresAt) {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "OTP không đúng hoặc đã hết hạn",
        }
    }
    delete(s.otpRepo, email)
    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Email đã được xác thực",
    }
}