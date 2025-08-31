package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/dto"
	"github.com/jackc/pgx/v5"
    "log"
	"github.com/jackc/pgx/v5/pgxpool"
)
type paymentRepository struct {
	db *pgxpool.Pool
}
func NewPaymentRepository(db *pgxpool.Pool) *paymentRepository {
	return &paymentRepository{db: db}
}
func (r *paymentRepository) CreatePayment(payment *dto.Payment) (dto.Payment, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `
        INSERT INTO payments (booking_id, amount, payment_method, status, transaction_id, paid_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING payment_id, booking_id, amount, payment_method, status, transaction_id, paid_at
    `

    var createdPayment dto.Payment
    err := r.db.QueryRow(ctx, query, payment.BookingID, payment.Amount, payment.PaymentMethod, payment.Status, payment.TransactionID, payment.PaidAt).Scan(
        &createdPayment.PaymentID,
        &createdPayment.BookingID,
        &createdPayment.Amount,
		&createdPayment.PaymentMethod,
        &createdPayment.Status,
        &createdPayment.TransactionID,
        &createdPayment.PaidAt,
    )
    if err != nil {
        return dto.Payment{}, fmt.Errorf("failed to create payment: %v", err)
    }

    log.Printf("Payment created successfully for booking_id %d", payment.BookingID)
    return createdPayment, nil

}

func (r *paymentRepository) UpdatePaymentStatus(paymentID int64, status string) (*dto.Payment, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    validStatuses := []string{"pending", "success", "failed"}
    isValid := false
    for _, validStatus := range validStatuses {
        if status == validStatus {
            isValid = true
            break
        }
    }
    if !isValid {
        return nil, fmt.Errorf("invalid status: %s", status)
    }

    query := `
    UPDATE payments
    SET status = $1, 
        paid_at = CASE 
                    WHEN $3 = 'success' THEN NOW() 
                    WHEN $3 = 'failed' THEN NULL 
                    ELSE paid_at 
                  END
    WHERE payment_id = $2
    RETURNING payment_id, booking_id, amount, payment_method, status, transaction_id, paid_at
`

    var updatedPayment dto.Payment
    err := r.db.QueryRow(ctx, query, status, paymentID, status).Scan(
        &updatedPayment.PaymentID,
        &updatedPayment.BookingID,
        &updatedPayment.Amount,
        &updatedPayment.PaymentMethod,
        &updatedPayment.Status,
        &updatedPayment.TransactionID,
        &updatedPayment.PaidAt,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, fmt.Errorf("payment not found: %v", err)
        }
        return nil, fmt.Errorf("failed to update payment status: %v", err)
    }

    log.Printf("Payment status updated successfully for payment_id %d", paymentID)
    return &updatedPayment, nil // Return pointer
}
func (r *paymentRepository) GetBookingIDByTransactionID(transactionID string) (int64, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `
        SELECT booking_id
        FROM payments
        WHERE transaction_id = $1
    `

    var bookingID int64
    err := r.db.QueryRow(ctx, query, transactionID).Scan(&bookingID)
    if err != nil {
        return 0, fmt.Errorf("failed to get booking ID for transaction ID %s: %v", transactionID, err)
    }

    return bookingID, nil
}
func (r *paymentRepository) GetPaymentIDByTransactionID(transactionID string) (int64, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    var paymentID int64
    query := `SELECT payment_id FROM payments WHERE transaction_id = $1`
    
    err := r.db.QueryRow(ctx, query, transactionID).Scan(&paymentID)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return 0, fmt.Errorf("payment not found for transaction ID %s", transactionID)
        }
        return 0, fmt.Errorf("failed to get payment ID by transaction ID: %w", err)
    }
    
    return paymentID, nil
}

func (r *paymentRepository) GetPaymentByBookingID(bookingID int64) (*dto.Payment, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `
        SELECT payment_id, booking_id, amount, payment_method, status, transaction_id, paid_at
        FROM payments
        WHERE booking_id = $1
        ORDER BY 
            CASE 
                WHEN status = 'success' THEN 1
                WHEN status = 'pending' THEN 2
                WHEN status = 'failed' THEN 3
                ELSE 4
            END,
            paid_at DESC
        LIMIT 1
    `

    var payment dto.Payment
    err := r.db.QueryRow(ctx, query, bookingID).Scan(
        &payment.PaymentID,
        &payment.BookingID,
        &payment.Amount,
        &payment.PaymentMethod,
        &payment.Status,
        &payment.TransactionID,
        &payment.PaidAt,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, fmt.Errorf("payment not found for booking ID %d", bookingID)
        }
        return nil, fmt.Errorf("failed to get payment by booking ID: %v", err)
    }

    return &payment, nil
}