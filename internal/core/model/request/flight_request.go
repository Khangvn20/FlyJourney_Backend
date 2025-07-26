package request
import (
	"time"
)

type CreateFlightRequest struct {
    AirlineID        int       `json:"airline_id" binding:"required"`
    FlightNumber     string    `json:"flight_number" binding:"required"`
    DepartureAirport string    `json:"departure_airport" binding:"required"`
    ArrivalAirport   string    `json:"arrival_airport" binding:"required"`
    DepartureTime    string    `json:"departure_time" binding:"required"`
    ArrivalTime      string   `json:"arrival_time" binding:"required"`
    DepartureCode    string    `json:"departure_airport_code" binding:"required"`
    ArrivalAirportCode string   `json:"arrival_airport_code" binding:"required"`
    DurationMinutes  int       `json:"duration_minutes" binding:"required"`
    StopsCount       int       `json:"stops_count"`
    TaxAndFees       float64   `json:"tax_and_fees"`
    Status           string    `json:"status" binding:"required"`
    Currency         string    `json:"currency" binding:"required"`
    Distance         int       `json:"distance"`
    FlightClasses    []FlightClassRequest `json:"flight_classes" binding:"required,dive"`
}
type FlightClassRequest struct {
    Class          string  `json:"class" binding:"required,oneof=economy premium_economy business first"`
    FareClassCode string `json:"fare_class_code" binding:"required"`
    BasePrice      float64 `json:"base_price" binding:"required"`
    AvailableSeats int     `json:"available_seats" binding:"required"`
    BasePriceInfant float64 `json:"base_price_infant" binding:"required"` // Optional, can be zero
    BasePriceChild float64 `json:"base_price_child" binding:"required"`
    TotalSeats     int     `json:"total_seats" binding:"required"`
    PackageAvailable string  `json:"package_available"`
}
type UpdateFlightRequest struct {
    AirlineID       int       `json:"airline_id" validate:"required"`
    AircraftID      int       `json:"aircraft_id" validate:"required"`
    FlightNumber    string    `json:"flight_number" validate:"required"`
    DepartureAirport string    `json:"departure_airport" validate:"required"`
    ArrivalAirport  string    `json:"arrival_airport" validate:"required"`
    DepartureTime   time.Time `json:"departure_time" validate:"required"`
    ArrivalTime     time.Time `json:"arrival_time" validate:"required"`
    DurationMinutes int       `json:"duration_minutes" validate:"required,min=1"`
    StopsCount      int       `json:"stops_count" validate:"required,min=0"`
    TotalSeats      int       `json:"total_seats" validate:"required,min=1"`
    TaxAndFees      float64   `json:"tax_and_fees" validate:"required,min=0"`
    Status          string    `json:"status" validate:"required,oneof=scheduled delayed cancelled boarding departed arrived diverted"`
    Distance        int       `json:"distance" validate:"required,min=1"`
}
type UpdateFlightStatusRequest struct {
    Status string `json:"status" binding:"required,oneof=scheduled delayed cancelled boarding departed arrived diverted"`
}

type FlightFilterRequest struct {
    AirlineIDs      []int    `json:"airline_ids"`
    MaxStops        int      `json:"max_stops"`
    MaxPrice        float64  `json:"max_price"`
    DepartureWindow []string `json:"departure_window"` // Format: ["06:00", "12:00"]
    ArrivalWindow   []string `json:"arrival_window"`   // Format: ["12:00", "18:00"]
}
type UpdateFlightClassRequest struct {
    Class          string  `json:"class" binding:"required,oneof=economy premium_economy business first"`
    BasePrice      float64 `json:"base_price" binding:"required"`
    AvailableSeats int     `json:"available_seats" binding:"required"`
    TotalSeats     int     `json:"total_seats" binding:"required"`
}
