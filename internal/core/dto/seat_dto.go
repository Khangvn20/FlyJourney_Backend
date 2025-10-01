package dto 

type SeatDTO struct {
	SeatID   int64     `json:"seat_id"`
	FlightClassID int64     `json:"flight_class_id"`
	SeatNumber string    `json:"seat_number"`
	Status	 string    `json:"status"`
}

