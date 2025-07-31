package request

type Passengers struct {
    Adults   int `json:"adults" binding:"required,min=1,max=9"`   
    Children int `json:"children" binding:"min=0,max=9"`       
    Infants  int `json:"infants" binding:"min=0,max=9"`     
}
type FlightSearchRequest struct {
    DepartureAirportCode string     `json:"departure_airport_code" binding:"required"`
    ArrivalAirportCode   string     `json:"arrival_airport_code" binding:"required"`
    DepartureDate    string `json:"departure_date" binding:"required"`
    ArrivalDate      string `json:"arrival_date" binding:"omitempty"`
    Passengers        Passengers   `json:"passenger" binding:"required"`
    FlightClass      string     `json:"flight_class" binding:"required"`   
    AirlineIDs       []int      `json:"airline_ids"`
    MaxStops         int        `json:"max_stops"`
    Page             int    `json:"page" binding:"omitempty,min=1"`
    Limit            int        `json:"limit" binding:"min=1,max=100"`
    SortBy           string     `json:"sort_by" binding:"omitempty,oneof=departure_time arrival_time price duration stops"`
    SortOrder string `json:"sort_order" binding:"omitempty,oneof=asc desc"`
}
type RoundtripFlightSearchRequest struct {
    DepartureAirportCode string     `json:"departure_airport_code" binding:"required"`
    ArrivalAirportCode   string     `json:"arrival_airport_code" binding:"required"`
    DepartureDate    string `json:"departure_date" binding:"required"`
    ReturnDate       string `json:"return_date" binding:"required"`
    FlightClass      string     `json:"flight_class" binding:"required"`
    AirlineIDs       []int      `json:"airline_ids"`
    MaxStops         int        `json:"max_stops"`
    Passengers       Passengers `json:"passengers" binding:"required"` // Changed from int to Passengers
    Page             int    `json:"page" binding:"omitempty,min=1"`
    Limit            int        `json:"limit" binding:"min=1,max=100"`
    SortBy           string     `json:"sort_by" binding:"omitempty,oneof=departure_time arrival_time price duration stops"`
    SortOrder string `json:"sort_order" binding:"omitempty,oneof=asc desc"`
}
type FlightCheapDaysRequest struct {
}
type FlightCheapMonthsRequest struct {
}