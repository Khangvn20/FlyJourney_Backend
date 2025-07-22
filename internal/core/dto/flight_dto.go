package dto
import "time"

type Flight struct {
    FlightID          int       `json:"flight_id"`
    AirlineID         int       `json:"airline_id"`
    FlightNumber      string    `json:"flight_number"`
    DepartureAirport  string    `json:"departure_airport"`
    ArrivalAiportCode string    `json:"arrival_airport_code"`
    DepartureAirportCode string `json:"departure_airport_code"`
    ArrivalAirport    string    `json:"arrival_airport"`
    DepartureTime     string    `json:"departure_time"`
    ArrivalTime       string    `json:"arrival_time"`
    DurationMinutes   int       `json:"duration_minutes"`
    StopsCount        int       `json:"stops_count"`
    TaxAndFees        float64   `json:"tax_and_fees"`
    TotalSeats        int       `json:"total_seats"`
    Status            string    `json:"status"`
    Currency          string    `json:"currency"`
    Distance          int       `json:"distance"`
    CreatedAt         time.Time `json:"created_at"`
    UpdatedAt         time.Time `json:"updated_at"`
}
type FlightClass struct {
    FlightClassID  int       `json:"flight_class_id"`
    FlightID       int       `json:"flight_id"`
    Class          string    `json:"class"` 
    BasePrice      float64   `json:"base_price"`
    BasePriceChild float64   `json:"base_price_child"`
    AvailableSeats int       `json:"available_seats"`
    PackageAvailable string  `json:"package_available,omitempty"`
    FreBaggage     string    `json:"free_baggage_allowance,omitempty"`
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
    DepartureTime     string `json:"departure_time"`
    ArrivalTime       string `json:"arrival_time"`
    DurationMinutes   int       `json:"duration_minutes"`
    StopsCount        int       `json:"stops_count"`
    BasePrice         float64   `json:"base_price"`
    TaxAndFees        float64   `json:"tax_and_fees"`
    TotalPrice        float64   `json:"total_price"`
    AvailableSeats    int       `json:"available_seats"`
    TotalFare         float64   `json:"total_fare"`
    Currency          string    `json:"currency"`
    TotalSeats        int       `json:"total_seats"`
    Status            string    `json:"status"`
    FlightClass       string    `json:"flight_class"`
    Distance          int       `json:"distance,omitempty"`
    ClassPrice        float64   `json:"class_price"`
    ClassAvailability int       `json:"class_availability"`
    PackageAvailable  string    `json:"package_available,omitempty"`

}
type RoundtripSearchResult struct {
    OutboundFlights []*FlightSearchResult `json:"outbound_flights"`
    InboundFlights  []*FlightSearchResult `json:"inbound_flights"`
}
