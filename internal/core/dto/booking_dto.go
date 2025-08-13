package dto

import (
	"time"
)

type Booking struct {
	BookingID      int64            `json:"booking_id"`
	UserID         int64            `json:"user_id"`
	FlightID       int64            `json:"flight_id"`
	BookingDate    time.Time        `json:"booking_date"`
	ContactEmail   string           `json:"contact_email"`    // Email liên hệ
	ContactPhone   string           `json:"contact_phone"`    // Số điện thoại liên hệ
	ContactAddress string           `json:"contact_address"`
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
	PassengerAge    int        `json:"passenger_age"`
	PassengerGender string     `json:"passenger_gender"`
	FlightClassID   int64      `json:"flight_class_id"`
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

type Payment struct {
	PaymentID     int64     `json:"payment_id"`
	BookingID     int64     `json:"booking_id"`
	Amount        float64   `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
	PaidAt        time.Time `json:"paid_at"`
	Status        string    `json:"status"`
	TransactionID string    `json:"transaction_id"`
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