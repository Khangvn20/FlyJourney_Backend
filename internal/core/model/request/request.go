package request

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"secret123"`
	Name     string `json:"name" binding:"required" example:"John Doe"`
	Phone    string `json:"phone,omitempty" example:"0123456789"`
}
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
