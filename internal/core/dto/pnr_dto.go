package dto

type PNR struct {
    PNRCode      string `json:"pnr_code"`
    BookingID    string `json:"booking_id"`
    FlightID     string `json:"flight_id"`
    ReturnFlightID string `json:"return_flight_id,omitempty"`
    Status       string `json:"status"`
    IssuedAt     string `json:"issued_at"`
    ExpiresAt    string `json:"expires_at"`
    CreatedBy    string `json:"created_by"`
    ModifiedBy   string `json:"modified_by,omitempty"`
    PNRData      string `json:"pnr_data,omitempty"`
}