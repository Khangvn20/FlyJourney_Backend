package request
type SendOTPRequest struct {
    Email string `json:"email" binding:"required,email"`
}

type VerifyOTPRequest struct {
    Email string `json:"email" binding:"required,email"`
    OTP   string `json:"otp" binding:"required"`

}