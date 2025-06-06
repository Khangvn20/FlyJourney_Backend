package response

import "github.com/Khangvn20/FlyJourney_Backend/internal/core/entity/error_code"

// Response defines the standard format for API responses.
type Response struct {
	Data         interface{}          `json:"data"`
	Status       bool                 `json:"status"`
	ErrorCode    error_code.ErrorCode `json:"errorCode"`
	ErrorMessage string               `json:"errorMessage"`
}

// RegisterResponse is the specific payload when registering a user successfully.
type RegisterResponse struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}
type LoginResponse struct {
	UserID   int    `json:"user_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Token    string `json:"token"`
	ExpireAt int64  `json:"expire_at"`
}
