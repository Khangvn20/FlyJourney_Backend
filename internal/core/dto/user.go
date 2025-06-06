package dto

import "time"

// User represents the user data structure used for data transfer between layers
type User struct {
	UserID    int        `json:"user_id"`
	Email     string     `json:"email" validate:"required,email"`
	Password  string     `json:"password,omitempty" validate:"required,min=8"`
	Name      string     `json:"name" validate:"required"`
	Phone     string     `json:"phone,omitempty"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	LastLogin *time.Time `json:"last_login,omitempty"`
}
