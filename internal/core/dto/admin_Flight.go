package dto
import "time"
type AdminFlightDetail struct {
    FlightID          int                    `json:"flight_id"`
    AirlineID         int                    `json:"airline_id"`
    AircraftID        int                    `json:"aircraft_id"`
    FlightNumber      string                 `json:"flight_number"`
    DepartureAirport  string                 `json:"departure_airport"`
    ArrivalAirport    string                 `json:"arrival_airport"`
    DepartureTime     string           `json:"departure_time"`
    ArrivalTime       string              `json:"arrival_time"`
    DurationMinutes   int                    `json:"duration_minutes"`
    StopsCount        int                    `json:"stops_count"`
    TaxAndFees        float64                `json:"tax_and_fees"`
    TotalSeats        int                    `json:"total_seats"`
    Status            string                 `json:"status"`
    Gate              string                 `json:"gate"`
    Terminal          string                 `json:"terminal"`
    Distance          int                    `json:"distance"`
    CreatedAt         time.Time              `json:"created_at"`
    UpdatedAt         time.Time              `json:"updated_at"`
    FlightClasses     []*AdminFlightClass    `json:"flight_classes"`
}

type AdminFlightClass struct {
    FlightClassID    int       `json:"flight_class_id"`
    FlightID         int       `json:"flight_id"`
    Class            string    `json:"class"`
    BasePrice        float64   `json:"base_price"`
    AvailableSeats   int       `json:"available_seats"`
    TotalSeats       int       `json:"total_seats"`
    PackageAvailable string    `json:"package_available,omitempty"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
}