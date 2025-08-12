package request
type CreateBookingRequest struct {
	UserID      int64	 `json:"user_id"`
	FlightID    int64	 `json:"flight_id" binding:"required"`
	ContactEmail string `json:"contact_email" binding:"required,email"`
    ContactAddress string `json:"contact_address" binding:"required"`
	ContactPhone   string `json:"contact_phone" binding:"required"`
	Note           string `json:"note"`
	TotalPrice    float64  `json:"total_price" binding:"required,numeric"`
	Details		[]*BookingDetailRequest `json:"details" binding:"required,dive"`
	 Ancillaries    []*AncillaryRequest     `json:"ancillaries" binding:"omitempty,dive"`
}
type BookingDetailRequest struct {
	PassengerAge     int       `json:"passenger_age" binding:"required,numeric"`
	PassengerGender  string    `json:"passenger_gender" binding:"required,oneof=male female"`
	FlightClassID    int64     `json:"flight_class_id" binding:"required,numeric"`
	Price            float64   `json:"price" binding:"required,numeric"`
	LastName         string    `json:"last_name" binding:"required"`
	FirstName       string    `json:"first_name" binding:"required"`
	DateOfBirth     string `json:"date_of_birth" binding:"required"`
	IDType          string    `json:"id_type" binding:"required"`
	IDNumber       string    `json:"id_number" binding:"required"`
	ExpiryDate     string `json:"expiry_date" binding:"required"`
	IssuingCountry  string    `json:"issuing_country" binding:"required"`
	Nationality     string    `json:"nationality" binding:"required"`
}
type AncillaryRequest struct {
	
	Type        string  `json:"type" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required,numeric"`
	Price      float64 `json:"price" binding:"required,numeric"`
}
