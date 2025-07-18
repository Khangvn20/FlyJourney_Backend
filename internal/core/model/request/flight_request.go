package request
import (
	"time"
    "encoding/json"
    "fmt"
)
type CustomTime struct {
    time.Time
}
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return fmt.Errorf("invalid time format: %v", err)
	}
	t, err := time.Parse("02-01-2006 15:04", str)
	if err != nil {
		return fmt.Errorf("invalid time format: %v", err)
	}
	ct.Time = t
	return nil
}
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(ct.Time.Format("02-01-2006 15:04"))
}
type CreateFlightRequest struct {AirlineID        int       `json:"airline_id" binding:"required"`
    AircraftID       int       `json:"aircraft_id" binding:"required"`
    FlightNumber     string    `json:"flight_number" binding:"required"`
    DepartureAirport string    `json:"departure_airport" binding:"required"`
    ArrivalAirport   string    `json:"arrival_airport" binding:"required"`
    DepartureTime    time.Time `json:"departure_time" binding:"required"`
    ArrivalTime      time.Time `json:"arrival_time" binding:"required"`
    DurationMinutes  int       `json:"duration_minutes" binding:"required"`
    StopsCount       int       `json:"stops_count"`
    TaxAndFees       float64   `json:"tax_and_fees"`
    Status           string    `json:"status" binding:"required"`
    Gate             string    `json:"gate"`
    Terminal         string    `json:"terminal"`
    Distance         int       `json:"distance"`
    FlightClasses    []FlightClassRequest `json:"flight_classes" binding:"required,dive"`
}
type FlightClassRequest struct {
    Class          string  `json:"class" binding:"required,oneof=economy premium_economy business first"`
    BasePrice      float64 `json:"base_price" binding:"required"`
    AvailableSeats int     `json:"available_seats" binding:"required"`
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
    TaxAndFees      float64   `json:"tax_and_fees" validate:"required,min=0"`
    Status          string    `json:"status" validate:"required,oneof=scheduled delayed cancelled boarding departed arrived diverted"`
    Gate            string    `json:"gate" validate:"required"`
    Terminal        string    `json:"terminal" validate:"required"`
    Distance        int       `json:"distance" validate:"required,min=1"`
}
type UpdateFlightStatusRequest struct {
    Status string `json:"status" binding:"required,oneof=scheduled delayed cancelled boarding departed arrived diverted"`
}
type FlightSearchRequest struct {
    DepartureAirport string     `json:"departure_airport" binding:"required"`
    ArrivalAirport   string     `json:"arrival_airport" binding:"required"`
    DepartureDate   CustomTime `json:"departure_date" binding:"required"`
    ArrivalDate     CustomTime `json:"arrival_date" binding:"omitempty"`
    FlightClass      string     `json:"flight_class" binding:"required"`
    AirlineIDs       []int      `json:"airline_ids"`
    MaxStops         int        `json:"max_stops"`
    Page             int        `json:"page"`
    Limit            int        `json:"limit"`
    SortBy           string     `json:"sort_by"`
    SortOrder        string     `json:"sort_order"`
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
type RoundtripFlightSearchRequest struct {
    DepartureAirport string     `json:"departure_airport" binding:"required"`
    ArrivalAirport   string     `json:"arrival_airport" binding:"required"`
    DepartureDate    CustomTime `json:"departure_date" binding:"required"`
    ReturnDate       CustomTime `json:"return_date" binding:"required"`
    FlightClass      string     `json:"flight_class" binding:"required"`
    AirlineIDs       []int      `json:"airline_ids"`
    MaxStops         int        `json:"max_stops"`
    Passengers       int       `json:"passengers" binding:"required,min=1"`
    Page             int       `json:"page"`
    Limit            int       `json:"limit"`
    SortBy           string    `json:"sort_by"`
    SortOrder        string    `json:"sort_order"`
}
type ConvertedRoundtripRequest struct {
    DepartureAirport string    `json:"departure_airport"`
    ArrivalAirport   string    `json:"arrival_airport"`
    DepartureDate    time.Time `json:"departure_date"`
    ReturnDate       time.Time `json:"return_date"`
    FlightClass      string    `json:"flight_class"`
    AirlineIDs       []int     `json:"airline_ids"`
    MaxStops         int       `json:"max_stops"`
    Passengers       int       `json:"passengers"`
    Page             int       `json:"page"`
    Limit            int       `json:"limit"`
    SortBy           string    `json:"sort_by"`
    SortOrder        string    `json:"sort_order"`
}