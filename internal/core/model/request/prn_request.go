package request 

type PNRRequest struct {
    BookingID    string `json:"booking_id" binding:"required"`
    FlightID     string `json:"flight_id" binding:"required"`
    ReturnFlightID string `json:"return_flight_id,omitempty"`
    Status       string `json:"status,omitempty"`
    IssuedAt     string `json:"issued_at,omitempty"`
    ExpiresAt    string `json:"expires_at,omitempty"`
    CreatedBy    string `json:"created_by" binding:"required"`
    PNRData      string `json:"pnr_data,omitempty"`
}