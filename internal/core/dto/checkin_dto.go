package dto
import "time"
type CheckinDTO struct {
	CheckinID  int64   `json:"checkin_id"`
	BookingDetailId int64   `json:"booking_detail_id"`
	FlightID       	int64   `json:"flight_id"`
	SeatID			int64   `json:"seat_id"`
	CheckinTIme    time.Time `json:"checkin_time"`
	CheckinMethod string    `json:"checkin_method"`
	Status        	string  `json:"status"`
	BoardingPassCode string  `json:"boarding_pass_code"`
	CreatedAt	 time.Time `json:"created_at"`
	UpdateAt	 time.Time `json:"update_at"`
}

type CheckinValidationResponse struct {
    IsEligible       bool                    `json:"is_eligible"`
    ErrorMessage     string                  `json:"error_message,omitempty"`
    PNRCode          string                  `json:"pnr_code"`
    FlightID         int64                   `json:"flight_id"`
    FlightNumber     string                  `json:"flight_number"`
    DepartureTime    time.Time               `json:"departure_time"`
    ArrivalTime      time.Time               `json:"arrival_time"`
    DepartureAirport string                  `json:"departure_airport"`
    ArrivalAirport   string                  `json:"arrival_airport"`
    AirlineName      string                  `json:"airline_name"`
    CheckinOpenTime  *time.Time              `json:"checkin_open_time,omitempty"`
    CheckinCloseTime *time.Time              `json:"checkin_close_time,omitempty"`
    BookingDetails   []*CheckinBookingDetail `json:"booking_details"`
}

// Booking detail for checkin
type CheckinBookingDetail struct {
    BookingDetailID     int64   `json:"booking_detail_id"`
    PassengerName       string  `json:"passenger_name"`
    PassengerAge        int     `json:"passenger_age"`
    PassengerGender     string  `json:"passenger_gender"`
    FlightClassName     string  `json:"flight_class_name"`
    SeatNumber          *string `json:"seat_number,omitempty"`
    IsCheckedIn         bool    `json:"is_checked_in"`
    CheckinID           *int64  `json:"checkin_id,omitempty"`
    BoardingPassCode    *string `json:"boarding_pass_code,omitempty"`
    IDNumber            string  `json:"id_number"`
    IDType              string  `json:"id_type"`
}

//Check seat map 
type ConfirmedSeatInfo struct {
    SeatNumber     string `json:"seat_number"`
    FlightClassID  int64  `json:"flight_class_id"`
    Status         string `json:"status"`
}

type SeatCheckResponse struct {
    FlightID        int64               `json:"flight_id"`
    FlightNumber    string              `json:"flight_number"`
    ConfirmedSeats  []ConfirmedSeatInfo `json:"confirmed_seats"`
}