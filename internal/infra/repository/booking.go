package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/dto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)
type bookingRepository struct {
	db *pgxpool.Pool
}
func NewBookingRepository(db *pgxpool.Pool) *bookingRepository {
	return &bookingRepository{db: db}
}
func (r *bookingRepository) CheckFlightClassAvailability(flightClassID int64) (bool, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
        SELECT available_seats, total_seats 
        FROM flight_classes 
        WHERE flight_class_id = $1
    `
	var available_seats, total_seats int
	err := r.db.QueryRow(ctx, query, flightClassID).Scan(&available_seats, &total_seats)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, 0, nil // No such flight class
		}
		return false, 0, fmt.Errorf("failed to check availability: %w", err)
	}
	return available_seats > 0, available_seats, nil
}

func (r *bookingRepository) CreateBooking(booking *dto.Booking) (*dto.Booking, error) {	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err :=r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	now :=time.Now()
	if booking.BookingDate.IsZero() {
		booking.BookingDate = now
	}
	booking.CreatedAt = now
	booking.UpdatedAt = now
	if booking.Status ==""{
		booking.Status = "pending"
	}
	if booking.CheckInStatus == "" {
		booking.CheckInStatus = "not_checked_in"
	}
    for _, detail := range booking.Details {
		var availableSeats int
		lockQuery := `SELECT available_seats FROM flight_classes 
                     WHERE flight_class_id = $1 FOR UPDATE`
        
        err := tx.QueryRow(ctx, lockQuery, detail.FlightClassID).Scan(&availableSeats)
        if err != nil {
            if errors.Is(err, pgx.ErrNoRows) {
                return nil, fmt.Errorf("flight class with ID %d not found", detail.FlightClassID)
            }
            return nil, fmt.Errorf("error locking flight class: %w", err)
        }
        
        if availableSeats < 1 {
            return nil, fmt.Errorf("no available seats for flight class %d", detail.FlightClassID)
        }
    }
	 bookingQuery := `
        INSERT INTO bookings (
            user_id, flight_id, booking_date, contact_email, 
            contact_phone, contact_address, note, status, 
            total_price, created_at, updated_at, check_in_status
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        RETURNING booking_id`

    err = tx.QueryRow(ctx, bookingQuery,
        booking.UserID, booking.FlightID, booking.BookingDate,
        booking.ContactEmail, booking.ContactPhone, booking.ContactAddress,
        booking.Note, booking.Status, booking.TotalPrice,
        booking.CreatedAt, booking.UpdatedAt, booking.CheckInStatus).Scan(&booking.BookingID)
    
    if err != nil {
        return nil, fmt.Errorf("error creating booking: %w", err)
    }
	 for i, detail := range booking.Details {
        detail.BookingID = booking.BookingID 

        detailQuery := `
            INSERT INTO booking_details (
                booking_id, passenger_age, passenger_gender, flight_class_id, 
                price, last_name, first_name, date_of_birth, id_type, 
                id_number, expiry_date, issuing_country, nationality
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
            RETURNING booking_detail_id`

        err = tx.QueryRow(ctx, detailQuery,
            detail.BookingID, detail.PassengerAge, detail.PassengerGender,
            detail.FlightClassID, detail.Price, detail.LastName, detail.FirstName,
            detail.DateOfBirth, detail.IDType, detail.IDNumber, detail.ExpiryDate,
            detail.IssuingCountry, detail.Nationality).Scan(&detail.BookingDetailID)
        
        if err != nil {
            return nil, fmt.Errorf("error creating booking detail: %w", err)
        }

		//Update seat availiable 
		 updateSeatsQuery := `
            UPDATE flight_classes 
            SET available_seats = available_seats - 1 
            WHERE flight_class_id = $1`

        _, err = tx.Exec(ctx, updateSeatsQuery, detail.FlightClassID)
        if err != nil {
            return nil, fmt.Errorf("error updating available seats: %w", err)
        }
		booking.Details[i] = detail // Cập nhật detail đã có booking_detail_id
	}
    if booking.Ancillaries != nil && len(booking.Ancillaries) > 0 { 
        for i, ancillary := range booking.Ancillaries {
            var bookingDetailID int64
        if len(booking.Details) > 0 {
            bookingDetailID = booking.Details[0].BookingDetailID
        } else {
            return nil, fmt.Errorf("no booking details to attach ancillary")
        }
            ancillary.CreatedAt = now
        ancillaryQuery := `
                INSERT INTO booking_ancillaries (
                    booking_detail_id, type, description, quantity, 
                    price, created_at
                ) VALUES ($1, $2, $3, $4, $5, $6)
                RETURNING ancillary_id`
        err = tx.QueryRow(ctx, ancillaryQuery,
            bookingDetailID, ancillary.Type, ancillary.Description,
            ancillary.Quantity, ancillary.Price, ancillary.CreatedAt).Scan(&ancillary.AncillaryID)
        if err != nil {
            return nil, fmt.Errorf("error creating ancillary: %w", err)
        }
    booking.TotalPrice += ancillary.Price * float64(ancillary.Quantity)
    //update ancillary with new ID
    ancillary.BookingDetailID = bookingDetailID
    booking.Ancillaries[i] = ancillary
	}
    //update total price
    updatePriceQuery := `UPDATE bookings SET total_price = $1 WHERE booking_id = $2`
        _, err = tx.Exec(ctx, updatePriceQuery, booking.TotalPrice, booking.BookingID)
        
        if err != nil {
            return nil, fmt.Errorf("error updating booking total price: %w", err)
        }
    }
	if err = tx.Commit(ctx); err != nil {
        return nil, fmt.Errorf("error committing transaction: %w", err)
    }
	return booking, nil
   
}