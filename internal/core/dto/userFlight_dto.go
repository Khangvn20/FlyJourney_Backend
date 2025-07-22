package dto

type UserFlightDetail struct {
    FlightID          int                    `json:"flight_id"`
    AirlineID         int                    `json:"airline_id"`
    FlightNumber      string                 `json:"flight_number"`
    DepartureAirport  string                 `json:"departure_airport"`
    ArrivalAirport    string                 `json:"arrival_airport"`
    DepartureTime     string           `json:"departure_time"`
    ArrivalTime       string             `json:"arrival_time"`
    DurationMinutes   int                    `json:"duration_minutes"`
    StopsCount        int                    `json:"stops_count"`
    TaxAndFees        float64                `json:"tax_and_fees"`
    Distance          int                    `json:"distance,omitempty"`
    FlightClasses     []*UserFlightClass     `json:"flight_classes"`
}
type UserFlightClass struct {
    FlightClassID    int     `json:"flight_class_id"`
    Class            string  `json:"class"`
    BasePrice        float64 `json:"base_price"`
    AvailableSeats   int     `json:"available_seats"`
    PackageAvailable string  `json:"package_available,omitempty"`
}