package dto

type TokenClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}
