package service
import "github.com/Khangvn20/FlyJourney_Backend/internal/core/model/response"

type EmailOTPService interface {
    SendOTPEmail(email string) *response.Response
    VerifyEmail(email, otp string) *response.Response
}