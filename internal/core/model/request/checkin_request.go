package request

type ValidateCheckin struct {
	PNRCode  string `json:"pnr_code" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required"`
}

