package request

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"secret123"`
	Name     string `json:"name" binding:"required" example:"John Doe"`
}
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type ConfirmRegisterRequest struct {
    Email    string `json:"email" binding:"required,email"`
    OTP      string `json:"otp" binding:"required"`
    Password string `json:"password" binding:"required,min=6"`
    Name     string `json:"name" binding:"required"`
}