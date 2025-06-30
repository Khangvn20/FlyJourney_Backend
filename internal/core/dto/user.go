package dto

import "time"
type UserRole string
const (
    RoleUser  UserRole = "user"
    RoleAdmin UserRole = "admin"
)

type User struct {
	UserID    int        `json:"user_id"`
	Email     string     `json:"email" validate:"required,email"`
	Password  string     `json:"password,omitempty" validate:"required,min=8"`
	Name      string     `json:"name" validate:"required"`
	Phone     string     `json:"phone,omitempty"`
	Roles      UserRole   `json:"roles"`  
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	
}
