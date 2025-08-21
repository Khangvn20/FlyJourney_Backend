package dto
import (
    "time"
    "database/sql"
)
type PNR struct {
    PNRCode      string `json:"pnr_code"`
    BookingID    string `json:"booking_id"`
    FlightID     string `json:"flight_id"`
    ReturnFlightID sql.NullString `json:"return_flight_id,omitempty"`
    Status       string `json:"status"`
    IssuedAt     *time.Time `json:"issued_at"`
    ExpiresAt    *time.Time `json:"expires_at"`
    CreatedBy    *time.Time `json:"created_by"`
    ModifiedBy   *time.Time `json:"modified_by,omitempty"`
    PNRData      *string `json:"pnr_data,omitempty"`
}