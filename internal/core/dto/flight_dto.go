package dto
import "time"

type Flight struct {
    FlightID          int       `json:"flight_id"`
    AirlineID         int       `json:"airline_id"`
	AircraftID        int       `json:"aircraft_id"`
    FlightNumber      string    `json:"flight_number"`
    DepartureAirport  string    `json:"departure_airport"`
    ArrivalAirport    string    `json:"arrival_airport"`
    DepartureTime     time.Time `json:"departure_time"`
    ArrivalTime       time.Time `json:"arrival_time"`
    DurationMinutes   int       `json:"duration_minutes"`
    StopsCount        int       `json:"stops_count"`
    TaxAndFees        float64   `json:"tax_and_fees"`
    TotalSeats        int       `json:"total_seats"`
    Status            string    `json:"status"`
    Gate              string    `json:"gate"`
    Terminal          string    `json:"terminal"`
    Distance          int       `json:"distance"`
    CreatedAt         time.Time `json:"created_at"`
    UpdatedAt         time.Time `json:"updated_at"`

}
type FlightClass struct {
    FlightClassID  int       `json:"flight_class_id"`
    FlightID       int       `json:"flight_id"`
    Class          string    `json:"class"` 
    BasePrice      float64   `json:"base_price"`
    AvailableSeats int       `json:"available_seats"`
    TotalSeats     int       `json:"total_seats"`
    CreatedAt      time.Time `json:"created_at,omitempty"`
    UpdatedAt      time.Time `json:"updated_at,omitempty"`
}
type FlightSearchResult struct {
    FlightID          int       `json:"flight_id"`
    AirlineID         int       `json:"airline_id"`
    AirlineName       string    `json:"airline_name"`
    FlightNumber      string    `json:"flight_number"`
    DepartureAirport  string    `json:"departure_airport"`
    ArrivalAirport    string    `json:"arrival_airport"`
    DepartureTime     time.Time `json:"departure_time"`
    ArrivalTime       time.Time `json:"arrival_time"`
    DurationMinutes   int       `json:"duration_minutes"`
    StopsCount        int       `json:"stops_count"`
    BasePrice         float64   `json:"base_price"`
    TaxAndFees        float64   `json:"tax_and_fees"`
    TotalPrice        float64   `json:"total_price"`
    AvailableSeats    int       `json:"available_seats"`
    TotalSeats        int       `json:"total_seats"`
    Status            string    `json:"status"`
    FlightClass       string    `json:"flight_class"`
    Gate              string    `json:"gate,omitempty"`
    Terminal          string    `json:"terminal,omitempty"`
    Distance          int       `json:"distance,omitempty"`
    ClassPrice        float64   `json:"class_price"`
    ClassAvailability int       `json:"class_availability"`
}