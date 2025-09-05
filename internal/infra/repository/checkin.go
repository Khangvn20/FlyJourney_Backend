package repository

import (
	"context"
	"errors"
	"fmt"
    "log"
	"time"

	"github.com/Khangvn20/FlyJourney_Backend/internal/core/dto"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)


type checkinRepository struct {
	db *pgxpool.Pool
}

func NewCheckinRepository(db *pgxpool.Pool) *checkinRepository {
    return &checkinRepository{db: db}
}
// ValidateCheckinByPNR - Validate checkin eligibility by PNR code, email, and name

func (r *checkinRepository) GetPNRInfo(pnrCode string) (*dto.PNRInfo, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `
        SELECT p.pnr_id, p.booking_id, p.flight_id, p.return_flight_id, p.status,
               b.contact_email, b.contact_name, b.status as booking_status,
               f.flight_number, f.departure_time, f.arrival_time, f.status,
               f.departure_airport, f.arrival_airport, f.departure_airport_code, f.arrival_airport_code,
               a.name as airline_name
        FROM pnr_records p
        JOIN bookings b ON p.booking_id = b.booking_id
        JOIN flights f ON p.flight_id = f.flight_id
        JOIN airlines a ON f.airline_id = a.airline_id
        WHERE p.pnr_code = $1
    `
    log.Printf("Executing query: %s with PNRCode: %s", query, pnrCode)
    
     var pnrInfo dto.PNRInfo
    err := r.db.QueryRow(ctx, query, pnrCode).Scan(
        &pnrInfo.PNRID, &pnrInfo.BookingID, &pnrInfo.FlightID, &pnrInfo.ReturnFlightID, &pnrInfo.Status,
        &pnrInfo.ContactEmail, &pnrInfo.ContactName, &pnrInfo.BookingStatus,
        &pnrInfo.FlightNumber, &pnrInfo.DepartureTime, &pnrInfo.ArrivalTime, &pnrInfo.FlightStatus,
        &pnrInfo.DepartureAirport, &pnrInfo.ArrivalAirport, &pnrInfo.DepartureCode, &pnrInfo.ArrivalCode,
        &pnrInfo.AirlineName,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, fmt.Errorf("PNR not found")
        }
        log.Printf("Error retrieving PNR info: %v", err)
        return nil, fmt.Errorf("error getting PNR info: %w", err)
    }
    log.Printf("PNR info retrieved: %+v", pnrInfo)
    return &pnrInfo, nil
}

// GetBookingDetailsForCheckin - CHỈ lấy booking details
func (r *checkinRepository) GetBookingDetailsForCheckin(bookingID int64) ([]*dto.CheckinBookingDetail, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `
        SELECT 
            bd.booking_detail_id, bd.first_name, bd.last_name, bd.passenger_age, 
            bd.passenger_gender, bd.id_number, bd.id_type,
            fc.class as flight_class_name,
            s.seat_number
        FROM booking_details bd
        JOIN flight_classes fc ON bd.flight_class_id = fc.flight_class_id
        LEFT JOIN seats s ON bd.seat_id = s.seat_id
        WHERE bd.booking_id = $1
        ORDER BY bd.booking_detail_id
    `

    log.Printf("Executing query: %s with BookingID: %d", query, bookingID)
    rows, err := r.db.Query(ctx, query, bookingID)
    if err != nil {
        return nil, fmt.Errorf("error querying booking details: %w", err)
    }
    defer rows.Close()

    var details []*dto.CheckinBookingDetail
    for rows.Next() {
        var detail dto.CheckinBookingDetail
        var firstName, lastName string
        var seatNumber *string

        err := rows.Scan(
            &detail.BookingDetailID, &firstName, &lastName, &detail.PassengerAge,
            &detail.PassengerGender, &detail.IDNumber, &detail.IDType,
            &detail.FlightClassName,
            &seatNumber,
        )
        if err != nil {
            log.Printf("Error scanning row for BookingID %d: %v", bookingID, err)
            continue
        }

        // Combine data
        detail.PassengerName = fmt.Sprintf("%s %s", firstName, lastName)
        detail.SeatNumber = seatNumber

        log.Printf("Retrieved booking detail: %+v", detail)
        details = append(details, &detail)
    }

    if err := rows.Err(); err != nil {
        log.Printf("Error iterating rows for BookingID %d: %v", bookingID, err)
        return nil, fmt.Errorf("error iterating rows: %w", err)
    }

    return details, nil
}