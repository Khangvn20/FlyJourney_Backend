package dto 
import "time"
type Booking struct {
	BookingID      int64     `json:"booking_id"`
	UserID         int64     `json:"user_id"`
	FlightID	   int64     `json:"flight_id"`
	BookingDate   time.Time  `json:"booking_date"`
	ContactEmail   string	`json:"contact_email"`      // Email liên hệ
	ContactPhone   string	`json:"contact_phone"`      // Số điện thoại liên hệ
	ContactAddress  string    `json:"contact_address"`  
	Note            string	`json:"note"`               // Ghi chú
	Status         string     `json:"status"`
	TotalPrice    float64    `json:"total_price"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Details      []*BookingDetail `json:"details"`
	Payment       *Payment `json:"payment"`
}

type BookingDetail struct {
    BookingDetailID int64     `json:"booking_detail_id"`
    BookingID       int64     `json:"booking_id"`
    PassengerAge    int       `json:"passenger_age"`
    PassengerGender string    `json:"passenger_gender"`
    FlightClassID   int64     `json:"flight_class_id"`
    SeatID          int64     `json:"seat_id"`
    Price           float64   `json:"price"`
    LastName        string    `json:"last_name"`
    FirstName       string    `json:"first_name"`
    DateOfBirth     string    `json:"date_of_birth"`    // ISO format: yyyy-mm-dd
    IDType          string    `json:"id_type"`
    IDNumber        string    `json:"id_number"`
    ExpiryDate      string    `json:"expiry_date"`      // ISO format: yyyy-mm-dd
    IssuingCountry  string    `json:"issuing_country"`
    Nationality     string    `json:"nationality"`
}
type Payment struct {
    PaymentID      int64     `json:"payment_id"`
    BookingID      int64     `json:"booking_id"`
    Amount        float64    `json:"amount"`
    PaymentMethod  string    `json:"payment_method"`
    PaidAt		time.Time	 `json:"paid_at"`
    Status         string    `json:"status"`
	TransactionID  string    `json:"transaction_id"`
}