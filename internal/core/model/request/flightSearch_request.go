package request


type FlightSearchRequest struct {
    DepartureAirport string     `json:"departure_airport" binding:"required"`
    ArrivalAirport   string     `json:"arrival_airport" binding:"required"`
    DepartureDate    string `json:"departure_date" binding:"required"`
    ArrivalDate      string `json:"arrival_date" binding:"omitempty"`
    FlightClass      string     `json:"flight_class" binding:"required"`   
    AirlineIDs       []int      `json:"airline_ids"`
    MaxStops         int        `json:"max_stops"`
    Page             int        `json:"page"`
    Limit            int        `json:"limit"`
    SortBy           string     `json:"sort_by"`
    SortOrder        string     `json:"sort_order"`
}
type RoundtripFlightSearchRequest struct {
    DepartureAirport string     `json:"departure_airport" binding:"required"`
    ArrivalAirport   string     `json:"arrival_airport" binding:"required"`
    DepartureDate    string `json:"departure_date" binding:"required"`
    ReturnDate       string `json:"return_date" binding:"required"`
    FlightClass      string     `json:"flight_class" binding:"required"`
    AirlineIDs       []int      `json:"airline_ids"`
    MaxStops         int        `json:"max_stops"`
    Passengers       int       `json:"passengers" binding:"required,min=1"`
    Page             int       `json:"page"`
    Limit            int       `json:"limit"`
    SortBy           string    `json:"sort_by"`
    SortOrder        string    `json:"sort_order"`
}
type FlightCheapDaysRequest struct {
}
type FlightCheapMonthsRequest struct {
}