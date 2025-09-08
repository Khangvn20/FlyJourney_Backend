package request

type ValidateCheckin struct {
	PNRCode  string `json:"pnr_code" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required"`
}

type OnlineCheckinRequest struct {
    BookingID int64                     `json:"booking_id" binding:"required"`
    Checkins  []CheckinRequest          `json:"checkins" binding:"required"`
}

type CheckinRequest struct {
    BookingDetailID int64  `json:"booking_detail_id" binding:"required"`
    SeatNumber      string `json:"seat_number" binding:"required"` // Frontend truyền seat number thay vì seat_id
}
