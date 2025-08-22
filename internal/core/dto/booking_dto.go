package dto

import (
	"time"
)

type Booking struct {
	BookingID      int64            `json:"booking_id"`
	UserID         int64            `json:"user_id"`
	FlightID       int64            `json:"flight_id"`
	ReturnFlightID *int64           `json:"return_flight_id"`
	BookingDate    time.Time        `json:"booking_date"`
	ContactEmail   string           `json:"contact_email"`    // Email liên hệ
	ContactPhone   string           `json:"contact_phone"`    // Số điện thoại liên hệ
	ContactAddress string           `json:"contact_address"`
	ContactName    string           `json:"contact_name"`     // Tên liên hệ
	Note           string           `json:"note"`             // Ghi chú
	Status         string           `json:"status"`           // e.g., "pending", "booked", "confirmed", "cancelled"
	TotalPrice     float64          `json:"total_price"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	CheckInStatus  string           `json:"check_in_status"`  // New: e.g., "not_checked_in", "checked_in", "seats_assigned"
	Details        []*BookingDetail `json:"details"`
	Payment        *Payment         `json:"payment"`
	Ancillaries     []*Ancillary     `json:"ancillaries"`
}


type BookingDetail struct {
	BookingDetailID int64      `json:"booking_detail_id"`
	BookingID       int64      `json:"booking_id"`
	FlightClassName string     `json:"flight_class_name"` // e.g., "Economy", "Business"
	PassengerAge    int        `json:"passenger_age"`
	PassengerGender string     `json:"passenger_gender"`
	FlightClassID   int64      `json:"flight_class_id"`
	ReturnFlightClassID *int64  `json:"return_flight_class_id"`
	ReturnFlightClassName *string `json:"return_flight_class_name"` // Optional: nil if not applicable
	SeatID          *int64     `json:"seat_id"`             // Optional: nil if not assigned yet
	Price           float64    `json:"price"`
	LastName        string     `json:"last_name"`
	FirstName       string     `json:"first_name"`
	DateOfBirth     time.Time  `json:"date_of_birth"`       // Changed to time.Time (parse from yyyy-mm-dd)
	IDType          string     `json:"id_type"`
	IDNumber        string     `json:"id_number"`
	ExpiryDate      time.Time  `json:"expiry_date"`         // Changed to time.Time (parse from yyyy-mm-dd)
	IssuingCountry  string     `json:"issuing_country"`
	Nationality     string     `json:"nationality"`
	//SeatAssignedAt  *time.Time `json:"seat_assigned_at"`    // New: Timestamp when seat is assigned (nil if not yet)
}
type Ancillary struct {
	AncillaryID   int64     `json:"ancillary_id"`   // Unique ID for this add-on instance
	BookingDetailID  int64     `json:"booking_detail_id"`
	Type          string    `json:"type"`           // e.g., "baggage", "meal", "priority_boarding"
	Description   string    `json:"description"`    // e.g., "Extra 10kg checked baggage", "Vegetarian meal"
	Quantity      int       `json:"quantity"`       // e.g., 1 for one meal, 2 for extra bags
	Price         float64   `json:"price"`          // Price for this add-on (can be dynamic)
	CreatedAt     time.Time `json:"created_at"`     // When this add-on was selected
}
type BookingEmailData struct {
	BookingID      int64     `json:"booking_id"`
	PNRCode       string    `json:"pnr_code"`
	UserFullName string    `json:"user_full_name"`
	ContactEmail  string    `json:"contact_email"`
	ContactPhone  string    `json:"contact_phone"`
	ContactAddress string   `json:"contact_address"`
	TotalPrice     string  `json:"total_price"`
	PaymentDate    time.Time `json:"payment_date"`
	PaymentMethod  string    `json:"payment_method"`
	OutboundFlight  *BookingEmailFlight `json:"outbound_flight"`
	InboundFlight   *BookingEmailFlight `json:"inbound_flight"`
	Passengers      []*BookingEmailPassenger   `json:"passengers"`
}
type BookingEmailFlight struct {
    FlightNumber     string    `json:"flight_number"`
    AirlineName      string    `json:"airline_name"`
    DepartureAirport string    `json:"departure_airport"`
    ArrivalAirport   string    `json:"arrival_airport"`
    DepartureTime    time.Time `json:"departure_time"`
    ArrivalTime      time.Time `json:"arrival_time"`
    FlightClass      string    `json:"flight_class"`
}

type BookingEmailPassenger struct {
    FullName    string `json:"full_name"`
    Type        string `json:"type"` // adult, child, infant
    SeatNumber  string `json:"seat_number,omitempty"`
	FlightClass  string `json:"flight_class,omitempty"`
}
