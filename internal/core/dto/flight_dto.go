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
    DepartureTime     time.Time    `json:"departure_time"`
    ArrivalTime       time.Time    `json:"arrival_time"`
    DurationMinutes   int       `json:"duration_minutes"`
    StopsCount        int       `json:"stops_count"`
    TaxAndFees        float64   `json:"tax_and_fees"`
    TotalSeats        int       `json:"total_seats"`
    Status            string    `json:"status"`
    Currency          string    `json:"currency"`
    Distance          int       `json:"distance"`
    FlightClasses []*FlightClass `json:"flight_classes,omitempty"`
    CreatedAt         time.Time `json:"created_at"`
    UpdatedAt         time.Time `json:"updated_at"`
}
type FlightClass struct {
    FlightClassID  int       `json:"flight_class_id"`
    FlightID       int       `json:"flight_id"`
    Class          string    `json:"class"` 
    BasePrice      float64   `json:"base_price"`
    FareClassCode  string    `json:"fare_class_code"`
    BasePriceChild float64   `json:"base_price_child"`
    BasePriceInfant float64 `json:"base_price_infant"`
    AvailableSeats int       `json:"available_seats"`
    TotalSeats     int       `json:"total_seats"`
    CreatedAt      time.Time `json:"created_at,omitempty"`
    FareClassDetails  *FareClasses  `json:"fare_class_details,omitempty"`
    UpdatedAt      time.Time `json:"updated_at,omitempty"`
}
type FlightSearchResult struct {
    // Flight basic info
    FlightID             int    `json:"flight_id"`
    FlightNumber         string `json:"flight_number"`
    AirlineID           int    `json:"airline_id"`
    AirlineName         string `json:"airline_name"`

    LogoUrl             string `json:"logo_url"`
    
    // Airport & Time info
    DepartureAirportCode string `json:"departure_airport_code"`
    ArrivalAirportCode   string `json:"arrival_airport_code"`
    DepartureAirport     string `json:"departure_airport"`
    ArrivalAirport       string `json:"arrival_airport"`
    DepartureTime        time.Time `json:"departure_time"`
    ArrivalTime          time.Time `json:"arrival_time"`
    
    // Flight details
    Duration             int    `json:"duration_minutes"`
    StopsCount          int    `json:"stops_count"`
    Distance            int    `json:"distance"`
    FlightClass         string `json:"flight_class"`
    TotalSeats          int    `json:"total_seats"`
    
    // Fare class info
    FareClassDetails    *FareClasses `json:"fare_class_details,omitempty"`
    
    // Pricing (gộp lại thành struct riêng)
    Pricing             PricingDetails `json:"pricing"`
    TaxAndFees        float64 `json:"tax_and_fees"`
}


type PricingDetails struct {
    BasePrices struct {
        Adult   float64 `json:"adult"`
        Child   float64 `json:"child"`
        Infant  float64 `json:"infant"`
    } `json:"base_prices"`
    
    TotalPrices struct {
        Adult   float64 `json:"adult"`
        Child   float64 `json:"child"`
        Infant  float64 `json:"infant"`
    } `json:"total_prices"`
    
    Taxes struct {
        Adult   float64 `json:"adult,omitempty"`
        Child   float64 `json:"child,omitempty"`
        Infant  float64 `json:"infant,omitempty"`
    } `json:"taxes,omitempty"`
    
    GrandTotal  float64 `json:"grand_total"`
    Currency    string  `json:"currency"`
}
type RoundtripSearchResult struct {
    OutboundFlights []*FlightSearchResult `json:"outbound_flights"`
    InboundFlights  []*FlightSearchResult `json:"inbound_flights"`
}
type FareClasses struct {
    FareClassCode string `json:"fare_class_code"`
    CabinClass    string  `json:"cabin_class"`
    Refundable    bool   `json:"refundable"`
    Changeable    bool   `json:"changeable"`
    Baggage_kg    string    `json:"baggage_kg"`
    Description   string    `json:"description"`
    RefundChangePolicy string `json:"refund_change_policy"`
    
}