package repository

import (
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/dto"
)

type PaymentRepository interface {
	CreatePayment(payment *dto.Payment)( dto.Payment, error)
	UpdatePaymentStatus(paymentID int64, status string) (*dto.Payment, error)
    GetBookingIDByTransactionID(transactionID string) (int64, error)
	GetPaymentIDByTransactionID(transactionID string) (int64, error)
}