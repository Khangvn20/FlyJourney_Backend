package dto

import "time"

type Payment struct {
	PaymentID     int64     `json:"payment_id"`
	BookingID     int64     `json:"booking_id"`
	Amount        string    `json:"amount"`
	Status        string    `json:"status"`
	TransactionID string    `json:"transaction_id"`
	PaymentMethod string    `json:"payment_method"`
	PaidAt        *time.Time`json:"paid_at"`
}