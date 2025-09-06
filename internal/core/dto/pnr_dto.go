package dto
import (
    "time"
    "database/sql"
)
type PNR struct {
    PNRID      int64 `json:"pnr_id"`
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

type PNRInfo struct {
    PNRID             int64     `json:"pnr_id"`
    BookingID         int64     `json:"booking_id"`
    FlightID          int64     `json:"flight_id"`
    ReturnFlightID    *int64    `json:"return_flight_id"`
    Status            string    `json:"status"`
    ContactEmail      string    `json:"contact_email"`
    ContactName       string    `json:"contact_name"`
    BookingStatus     string    `json:"booking_status"`
    FlightNumber      string    `json:"flight_number"`
    DepartureTime     time.Time `json:"departure_time"`
    ArrivalTime       time.Time `json:"arrival_time"`
    FlightStatus      string    `json:"flight_status"`
    DepartureAirport  string    `json:"departure_airport"`
    ArrivalAirport    string    `json:"arrival_airport"`
    DepartureCode     string    `json:"departure_code"`
    ArrivalCode       string    `json:"arrival_code"`
    AirlineName       string    `json:"airline_name"`
}