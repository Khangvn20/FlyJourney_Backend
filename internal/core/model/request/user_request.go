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
type ConfirmResetPasswordRequest struct {
	Email    string `json:"email" binding:"required,email"`
	OTP      string `json:"otp" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
type ResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}
type UpdateProfileRequest struct {
    Name  string `json:"name" `
    Phone string `json:"phone" `
	Email string `json:"email"`
}