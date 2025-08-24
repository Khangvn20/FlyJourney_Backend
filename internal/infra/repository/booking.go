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
			return false, 0, nil 
		}
		return false, 0, fmt.Errorf("failed to check availability: %w", err)
	}
	return available_seats > 0, available_seats, nil
}

func (r *bookingRepository) CreateBooking(booking *dto.Booking) (*dto.Booking, error) {	
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
        if detail.ReturnFlightClassID != nil {
            err :=tx.QueryRow(ctx, lockQuery, *detail.ReturnFlightClassID).Scan(&availableSeats)
            if err != nil {
                if errors.Is(err, pgx.ErrNoRows) {
                    return nil, fmt.Errorf("return flight class with ID %d not found", *detail.ReturnFlightClassID)
                }
                return nil, fmt.Errorf("error locking return flight class: %w", err)
            }
            if availableSeats < 1 {
            return nil, fmt.Errorf("no available seats for return flight class %d", *detail.ReturnFlightClassID)
            }
        }
    }
	 bookingQuery := `
        INSERT INTO bookings (
            user_id, flight_id, return_flight_id, booking_date, contact_email, 
            contact_phone, contact_address, contact_name, note, status, 
            total_price, created_at, updated_at, check_in_status
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
        RETURNING booking_id`

    err = tx.QueryRow(ctx, bookingQuery,
        booking.UserID, booking.FlightID, booking.ReturnFlightID, booking.BookingDate,
        booking.ContactEmail, booking.ContactPhone, booking.ContactAddress, booking.ContactName,
        booking.Note, booking.Status, booking.TotalPrice,
        booking.CreatedAt, booking.UpdatedAt, booking.CheckInStatus).Scan(&booking.BookingID)
    
    if err != nil {
        return nil, fmt.Errorf("error creating booking: %w", err)
    }
	 for i, detail := range booking.Details {
        detail.BookingID = booking.BookingID 

        detailQuery := `
            INSERT INTO booking_details (
                booking_id, passenger_age, passenger_gender, flight_class_id,return_flight_class_id,
                price, last_name, first_name, date_of_birth, id_type, 
                id_number, expiry_date, issuing_country, nationality
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
            RETURNING booking_detail_id`

        err = tx.QueryRow(ctx, detailQuery,
            detail.BookingID, detail.PassengerAge, detail.PassengerGender,
            detail.FlightClassID, detail.ReturnFlightClassID, detail.Price, detail.LastName, detail.FirstName,
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
        if detail.ReturnFlightClassID != nil {
    if _, err = tx.Exec(ctx, updateSeatsQuery, *detail.ReturnFlightClassID); err != nil {
        return nil, fmt.Errorf("error updating available seats for return flight class %d: %w", *detail.ReturnFlightClassID, err)
    }
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
func (r *bookingRepository) GetBookingByID(bookingID int64) (*dto.Booking, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    query := `
        SELECT 
            booking_id, user_id, flight_id, return_flight_id, booking_date, 
            contact_email, contact_phone, contact_address, note, status, 
            total_price,  check_in_status
        FROM bookings
        WHERE booking_id = $1
    `

    var booking dto.Booking
    err := r.db.QueryRow(ctx, query, bookingID).Scan(
        &booking.BookingID,
        &booking.UserID,
        &booking.FlightID,
        &booking.ReturnFlightID,
        &booking.BookingDate,
        &booking.ContactEmail,
        &booking.ContactPhone,
        &booking.ContactAddress,
        &booking.Note,
        &booking.Status,
        &booking.TotalPrice,
        &booking.CheckInStatus,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, fmt.Errorf("booking with ID %d not found", bookingID)
        }
        return nil, fmt.Errorf("error getting booking by ID: %w", err)
    }
      detailsQuery := `
        SELECT 
            bd.booking_detail_id, bd.booking_id, bd.passenger_age, bd.passenger_gender,
            bd.flight_class_id, bd.return_flight_class_id, bd.price, bd.last_name, 
            bd.first_name, bd.date_of_birth, bd.id_type, bd.id_number, bd.expiry_date, 
            bd.issuing_country, bd.nationality,
            fc.class as flight_class_name,
            rfc.class as return_flight_class_name
        FROM booking_details bd
        LEFT JOIN flight_classes fc ON bd.flight_class_id = fc.flight_class_id
        LEFT JOIN flight_classes rfc ON bd.return_flight_class_id = rfc.flight_class_id
        WHERE bd.booking_id = $1
        ORDER BY bd.booking_detail_id
    `

    rows, err := r.db.Query(ctx, detailsQuery, bookingID)
    if err != nil {
        log.Printf("Error fetching booking details: %v", err)
        return nil, fmt.Errorf("error fetching booking details: %w", err)
    }
    defer rows.Close()

    var details []*dto.BookingDetail
    for rows.Next() {
        var detail dto.BookingDetail
        var flightClassName, returnFlightClassName *string

        err := rows.Scan(
            &detail.BookingDetailID,
            &detail.BookingID,
            &detail.PassengerAge,
            &detail.PassengerGender,
            &detail.FlightClassID,
            &detail.ReturnFlightClassID,
            &detail.Price,
            &detail.LastName,
            &detail.FirstName,
            &detail.DateOfBirth,
            &detail.IDType,
            &detail.IDNumber,
            &detail.ExpiryDate,
            &detail.IssuingCountry,
            &detail.Nationality,
            &flightClassName,
            &returnFlightClassName,
        )
        if err != nil {
            log.Printf("Error scanning booking detail: %v", err)
            continue
        }

        // Gán flight class name nếu có
        if flightClassName != nil {
            detail.FlightClassName = *flightClassName
        } else {
            detail.FlightClassName = "Unknown"
        }
        
        if returnFlightClassName != nil {
            detail.ReturnFlightClassName = returnFlightClassName
        }
        
        details = append(details, &detail)
        log.Printf("Loaded booking detail: ID=%d, Name=%s %s, Age=%d, FlightClass=%s (ID=%d)", 
            detail.BookingDetailID, detail.FirstName, detail.LastName, detail.PassengerAge, 
            detail.FlightClassName, detail.FlightClassID)
    }

    booking.Details = details
    log.Printf("Total booking details loaded for booking %d: %d", bookingID, len(booking.Details))

    if len(details) > 0 {
        ancillariesQuery := `
            SELECT 
                ancillary_id, booking_detail_id, type, description, 
                quantity, price, created_at
            FROM booking_ancillaries ba
            WHERE ba.booking_detail_id IN (
                SELECT booking_detail_id FROM booking_details WHERE booking_id = $1
            )
        `

        ancillaryRows, err := r.db.Query(ctx, ancillariesQuery, bookingID)
        if err != nil {
            log.Printf("Error fetching ancillaries: %v", err)
        } else {
            defer ancillaryRows.Close()
            
            var ancillaries []*dto.Ancillary
            for ancillaryRows.Next() {
                var ancillary dto.Ancillary
                err := ancillaryRows.Scan(
                    &ancillary.AncillaryID,
                    &ancillary.BookingDetailID,
                    &ancillary.Type,
                    &ancillary.Description,
                    &ancillary.Quantity,
                    &ancillary.Price,
                    &ancillary.CreatedAt,
                )
                if err != nil {
                    log.Printf("Error scanning ancillary: %v", err)
                    continue
                }
                ancillaries = append(ancillaries, &ancillary)
            }
            booking.Ancillaries = ancillaries
        }
    }

    return &booking, nil
}

func (r *bookingRepository) GetExpiredBookingIDs() ([]int64, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `
        SELECT booking_id
        FROM bookings
        WHERE status = 'pending_payment' AND created_at <= NOW() - INTERVAL '2 hours'
    `

    rows, err := r.db.Query(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("error fetching expired booking IDs: %w", err)
    }
    defer rows.Close()

    var expiredBookingIDs []int64
    for rows.Next() {
        var bookingID int64
        if err := rows.Scan(&bookingID); err != nil {
            return nil, fmt.Errorf("error scanning expired booking ID: %w", err)
        }
        expiredBookingIDs = append(expiredBookingIDs, bookingID)
    }

    return expiredBookingIDs, nil
}

func (r *bookingRepository) CancelBookings(bookingIDs []int64) ([]int64, error) {
    log.Printf("Canceling bookings with IDs: %v", bookingIDs)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    query := `
        DELETE FROM bookings
        WHERE booking_id = ANY($1) AND status = 'pending_payment'
        RETURNING booking_id
    `

    rows, err := r.db.Query(ctx, query, bookingIDs)
    if err != nil {
        return nil, fmt.Errorf("error canceling bookings: %w", err)
    }
    defer rows.Close()

    var canceledBookingIDs []int64
    for rows.Next() {
        var bookingID int64
        if err := rows.Scan(&bookingID); err != nil {
            return nil, fmt.Errorf("error scanning canceled booking ID: %w", err)
        }
        canceledBookingIDs = append(canceledBookingIDs, bookingID)
    }

    return canceledBookingIDs, nil
}
func(r *bookingRepository) UpdateStatusConfirm(bookingID int64) (*dto.Booking, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `
        UPDATE bookings
        SET status = 'confirmed', updated_at = NOW()
        WHERE booking_id = $1
        RETURNING booking_id, user_id, status, updated_at
    `

    var booking dto.Booking
    err := r.db.QueryRow(ctx, query, bookingID).Scan(&booking.BookingID, &booking.UserID, &booking.Status, &booking.UpdatedAt)
    if err != nil {
        return nil, fmt.Errorf("error updating booking status: %w", err)
    }
    
    return &booking, nil
}
func (r *bookingRepository) UpdateBookingStatus(bookingID int64, status string) (*dto.Booking, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `
        UPDATE bookings
        SET status = $1, updated_at = NOW()
        WHERE booking_id = $2
        RETURNING booking_id, user_id, status, updated_at
    `

    var booking dto.Booking
    err := r.db.QueryRow(ctx, query, status, bookingID).Scan(&booking.BookingID, &booking.UserID, &booking.Status, &booking.UpdatedAt)
    if err != nil {
        return nil, fmt.Errorf("error updating booking status: %w", err)
    }

    return &booking, nil
}
func (r *bookingRepository) GetBookingsByFlightID(flightID int64) ([]*dto.Booking, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `
        SELECT 
            booking_id, user_id, flight_id, return_flight_id, booking_date, 
            contact_email, contact_phone, contact_address, note, status, 
            total_price, check_in_status, created_at, updated_at
        FROM bookings
        WHERE flight_id = $1
        ORDER BY booking_date DESC
    `

    rows, err := r.db.Query(ctx, query, flightID)
    if err != nil {
        return nil, fmt.Errorf("error fetching bookings by flight ID: %w", err)
    }
    defer rows.Close()

    var bookings []*dto.Booking
    for rows.Next() {
        var booking dto.Booking
        err := rows.Scan(
            &booking.BookingID,
            &booking.UserID,
            &booking.FlightID,
            &booking.ReturnFlightID,
            &booking.BookingDate,
            &booking.ContactEmail,
            &booking.ContactPhone,
            &booking.ContactAddress,
            &booking.Note,
            &booking.Status,
            &booking.TotalPrice,
            &booking.CheckInStatus,
            &booking.CreatedAt,
            &booking.UpdatedAt,
        )
        if err != nil {
            log.Printf("Error scanning booking: %v", err)
            continue
        }

        bookings = append(bookings, &booking)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating over bookings: %w", err)
    }

    return bookings, nil
}
func (r *bookingRepository) GetAllBookingByUserID(userID int64) ([]*dto.Booking, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Query để lấy thông tin cơ bản của booking
    query := `
        SELECT 
            booking_id, user_id, flight_id, return_flight_id, booking_date, 
            contact_email, contact_phone, contact_address, note, status, 
            total_price, check_in_status, created_at, updated_at
        FROM bookings
        WHERE user_id = $1
        ORDER BY booking_date DESC
    `

    rows, err := r.db.Query(ctx, query, userID)
    if err != nil {
        return nil, fmt.Errorf("error fetching bookings by user ID: %w", err)
    }
    defer rows.Close()

    var bookings []*dto.Booking
    bookingMap := make(map[int64]*dto.Booking)

    for rows.Next() {
        var booking dto.Booking
        err := rows.Scan(
            &booking.BookingID,
            &booking.UserID,
            &booking.FlightID,
            &booking.ReturnFlightID,
            &booking.BookingDate,
            &booking.ContactEmail,
            &booking.ContactPhone,
            &booking.ContactAddress,
            &booking.Note,
            &booking.Status,
            &booking.TotalPrice,
            &booking.CheckInStatus,
            &booking.CreatedAt,
            &booking.UpdatedAt,
        )
        if err != nil {
            log.Printf("Error scanning booking: %v", err)
            continue
        }

        bookings = append(bookings, &booking)
        bookingMap[booking.BookingID] = &booking
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating over bookings: %w", err)
    }


    detailsQuery := `
        SELECT 
            bd.booking_detail_id, bd.booking_id, bd.passenger_age, bd.passenger_gender,
            bd.flight_class_id, bd.return_flight_class_id, bd.price, bd.last_name, 
            bd.first_name, bd.date_of_birth, bd.id_type, bd.id_number, bd.expiry_date, 
            bd.issuing_country, bd.nationality,
            fc.class as flight_class_name,
            rfc.class as return_flight_class_name
        FROM booking_details bd
        LEFT JOIN flight_classes fc ON bd.flight_class_id = fc.flight_class_id
        LEFT JOIN flight_classes rfc ON bd.return_flight_class_id = rfc.flight_class_id
        WHERE bd.booking_id = ANY($1)
        ORDER BY bd.booking_detail_id
    `

    bookingIDs := make([]int64, 0, len(bookingMap))
    for id := range bookingMap {
        bookingIDs = append(bookingIDs, id)
    }

    detailRows, err := r.db.Query(ctx, detailsQuery, bookingIDs)
    if err != nil {
        return nil, fmt.Errorf("error fetching booking details: %w", err)
    }
    defer detailRows.Close()

    for detailRows.Next() {
        var detail dto.BookingDetail
        var flightClassName, returnFlightClassName *string

        err := detailRows.Scan(
            &detail.BookingDetailID,
            &detail.BookingID,
            &detail.PassengerAge,
            &detail.PassengerGender,
            &detail.FlightClassID,
            &detail.ReturnFlightClassID,
            &detail.Price,
            &detail.LastName,
            &detail.FirstName,
            &detail.DateOfBirth,
            &detail.IDType,
            &detail.IDNumber,
            &detail.ExpiryDate,
            &detail.IssuingCountry,
            &detail.Nationality,
            &flightClassName,
            &returnFlightClassName,
        )
        if err != nil {
            log.Printf("Error scanning booking detail: %v", err)
            continue
        }

        if flightClassName != nil {
            detail.FlightClassName = *flightClassName
        }
        if returnFlightClassName != nil {
            detail.ReturnFlightClassName = returnFlightClassName
        }
        if booking, exists := bookingMap[detail.BookingID]; exists {
            booking.Details = append(booking.Details, &detail)
        }
    }
//query payment
paymentsQuery := `
        SELECT 
            payment_id, booking_id, amount, status, payment_method ,paid_at, transaction_id
        FROM payments
        WHERE booking_id = ANY($1)
    `

    bookingIDs = make([]int64, 0, len(bookingMap))
    for id := range bookingMap {
    bookingIDs = append(bookingIDs, id)
    }
    paymentRows, err := r.db.Query(ctx, paymentsQuery, bookingIDs)
    if err != nil {
        return nil, fmt.Errorf("error fetching payments: %w", err)
    }
    defer paymentRows.Close()

    for paymentRows.Next() {
        var payment dto.Payment
        err := paymentRows.Scan(
            &payment.PaymentID,
            &payment.BookingID,
            &payment.Amount,
            &payment.Status,
            &payment.PaymentMethod,
            &payment.PaidAt,
            &payment.TransactionID,
        )
        if err != nil {
            log.Printf("Error scanning payment: %v", err)
            continue
        }

        if booking, exists := bookingMap[payment.BookingID]; exists {
            booking.Payment = &payment
        }
    }

    ancillariesQuery := `
        SELECT 
            ancillary_id, booking_detail_id, type, description, 
            quantity, price, created_at
        FROM booking_ancillaries ba
        WHERE ba.booking_detail_id = ANY(
            SELECT booking_detail_id FROM booking_details WHERE booking_id = ANY($1)
        )
    `

    ancillaryRows, err := r.db.Query(ctx, ancillariesQuery, bookingIDs)
    if err != nil {
        return nil, fmt.Errorf("error fetching ancillaries: %w", err)
    }
    defer ancillaryRows.Close()

    for ancillaryRows.Next() {
        var ancillary dto.Ancillary
        err := ancillaryRows.Scan(
            &ancillary.AncillaryID,
            &ancillary.BookingDetailID,
            &ancillary.Type,
            &ancillary.Description,
            &ancillary.Quantity,
            &ancillary.Price,
            &ancillary.CreatedAt,
        )
        if err != nil {
            log.Printf("Error scanning ancillary: %v", err)
            continue
        }

        for _, booking := range bookings {
            for _, detail := range booking.Details {
                if detail.BookingDetailID == ancillary.BookingDetailID {
                    booking.Ancillaries = append(booking.Ancillaries, &ancillary)
                    break
                }
            }
        }
    }

    return bookings, nil
}