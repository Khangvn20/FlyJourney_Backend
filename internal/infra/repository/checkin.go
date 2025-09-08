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

func (r *checkinRepository) GetConfirmedSeatsByFlightID(flightID int64) (*dto.SeatCheckResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Query đơn giản để lấy thông tin flight
    flightQuery := `SELECT flight_id, flight_number FROM flights WHERE flight_id = $1`
    var flightIDResult int64
    var flightNumberResult string
    
    err := r.db.QueryRow(ctx, flightQuery, flightID).Scan(&flightIDResult, &flightNumberResult)
    if err != nil {
        if err == pgx.ErrNoRows {
            log.Printf("Flight not found for FlightID: %d", flightID)
            return nil, fmt.Errorf("flight not found")
        }
        log.Printf("Error getting flight info: %v", err)
        return nil, fmt.Errorf("error getting flight info: %w", err)
    }

    // Query đơn giản để lấy ghế confirmed từ bảng seats
    seatQuery := `
    SELECT s.seat_number, s.flight_class_id, s.status, fc.class, fc.fare_class_code
    FROM seats s
    JOIN flight_classes fc ON s.flight_class_id = fc.flight_class_id
    WHERE s.flight_id = $1 AND s.status = 'confirm'
    ORDER BY s.seat_number
`

    log.Printf("Executing confirmed seats query for FlightID: %d", flightID)
    
    rows, err := r.db.Query(ctx, seatQuery, flightID)
    if err != nil {
        log.Printf("Error executing confirmed seats query: %v", err)
        return nil, fmt.Errorf("error querying confirmed seats: %w", err)
    }
    defer rows.Close()

    var confirmedSeats []dto.ConfirmedSeatInfo

    for rows.Next() {
    var seatNumber string
    var flightClassID int64
    var status string
    var class string
    var fareClassCode string

    err := rows.Scan(&seatNumber, &flightClassID, &status, &class, &fareClassCode)
    if err != nil {
        log.Printf("Error scanning confirmed seat row: %v", err)
        continue
    }

    confirmedSeat := dto.ConfirmedSeatInfo{
        SeatNumber:    seatNumber,
        FlightClassID: flightClassID,
        Status:        status,
        Class:         class,         
        FareClassCode: fareClassCode, 
    }
        confirmedSeats = append(confirmedSeats, confirmedSeat)
    }

    if err := rows.Err(); err != nil {
        log.Printf("Error iterating rows for FlightID %d: %v", flightID, err)
        return nil, fmt.Errorf("error iterating rows: %w", err)
    }

    response := &dto.SeatCheckResponse{
        FlightID:       flightIDResult,
        FlightNumber:   flightNumberResult,
        ConfirmedSeats: confirmedSeats,
    }

    log.Printf("Retrieved %d confirmed seats for FlightID: %d", len(confirmedSeats), flightID)
    
    return response, nil
}
func (r *checkinRepository) CheckSeatOccupied(flightID int64, seatNumber string) (bool, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `
        SELECT EXISTS(
            SELECT 1 FROM seats 
            WHERE flight_id = $1 AND seat_number = $2 AND status = 'confirm'
        )
    `
    
    var occupied bool
    err := r.db.QueryRow(ctx, query, flightID, seatNumber).Scan(&occupied)
    if err != nil {
        return false, fmt.Errorf("error checking seat occupancy: %w", err)
    }
    
    return occupied, nil
}

// CheckAlreadyCheckedIn - Kiểm tra đã check-in chưa
func (r *checkinRepository) CheckAlreadyCheckedIn(bookingDetailID int64) (bool, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `SELECT EXISTS(SELECT 1 FROM check_ins WHERE booking_detail_id = $1)`
    
    var checkedIn bool
    err := r.db.QueryRow(ctx, query, bookingDetailID).Scan(&checkedIn)
    return checkedIn, err
}

// AssignSeatToBookingDetail - Gán ghế cho passenger và tạo record trong seats table
func (r *checkinRepository) AssignSeatToBookingDetail(bookingDetailID int64, seatNumber string, flightID int64) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    tx, err := r.db.Begin(ctx)
    if err != nil {
        return fmt.Errorf("error starting transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    // 1. Lấy flight_class_id từ booking_detail
    var flightClassID int64
    flightClassQuery := `SELECT flight_class_id FROM booking_details WHERE booking_detail_id = $1`
    err = tx.QueryRow(ctx, flightClassQuery, bookingDetailID).Scan(&flightClassID)
    if err != nil {
        return fmt.Errorf("error getting flight class: %w", err)
    }

    // 2. Insert seat record vào database (tạo mới)
    insertSeatQuery := `
        INSERT INTO seats (flight_id, flight_class_id, seat_number, status)
        VALUES ($1, $2, $3, 'confirm')
        RETURNING seat_id
    `
    var seatID int64
    err = tx.QueryRow(ctx, insertSeatQuery, flightID, flightClassID, seatNumber).Scan(&seatID)
    if err != nil {
        return fmt.Errorf("error creating seat record: %w", err)
    }

    // 3. Update booking_details với seat_id
    updateBookingQuery := `
        UPDATE booking_details 
        SET seat_id = $1
        WHERE booking_detail_id = $2
    `
    _, err = tx.Exec(ctx, updateBookingQuery, seatID, bookingDetailID)
    if err != nil {
        return fmt.Errorf("error updating booking detail: %w", err)
    }

    return tx.Commit(ctx)
}

// CreateCheckinRecord - Tạo record check-in
func (r *checkinRepository) CreateCheckinRecord(checkinData *dto.CheckinDTO) (*dto.CheckinDTO, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
 
    query := `
        INSERT INTO check_ins (
            booking_detail_id, seat_id, flight_id, status, check_in_time, 
            boarding_pass_code, check_in_method, created_at, updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING check_in_id
    `
    
    err := r.db.QueryRow(ctx, query,
        checkinData.BookingDetailId, checkinData.SeatID, checkinData.FlightID, checkinData.Status,
        checkinData.CheckinTIme, checkinData.BoardingPassCode, checkinData.CheckinMethod,
        checkinData.CreatedAt, checkinData.UpdateAt,
    ).Scan(&checkinData.CheckinID)
    
    return checkinData, err
}

// GetFlightInfoByBooking - Lấy thông tin flight từ booking
func (r *checkinRepository) GetFlightInfoByBooking(bookingID int64) (*dto.FlightBasicInfo, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `
        SELECT b.flight_id, f.flight_number, f.departure_time, b.status
        FROM bookings b
        JOIN flights f ON b.flight_id = f.flight_id
        WHERE b.booking_id = $1
    `
    
    var flight dto.FlightBasicInfo
    err := r.db.QueryRow(ctx, query, bookingID).Scan(
        &flight.FlightID, &flight.FlightNumber, 
        &flight.DepartureTime, &flight.BookingStatus,
    )
    
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, fmt.Errorf("booking not found")
        }
        return nil, fmt.Errorf("error getting flight info: %w", err)
    }
    
    return &flight, nil
}

// ValidateBookingDetailOwnership - Validate booking detail thuộc về booking
func (r *checkinRepository) ValidateBookingDetailOwnership(bookingDetailID int64, bookingID int64) (*dto.PassengerInfo, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `
        SELECT bd.booking_detail_id, bd.first_name, bd.last_name, fc.class
        FROM booking_details bd
        JOIN flight_classes fc ON bd.flight_class_id = fc.flight_class_id
        WHERE bd.booking_detail_id = $1 AND bd.booking_id = $2
    `
    
    var passenger dto.PassengerInfo
    var firstName, lastName string
    
    err := r.db.QueryRow(ctx, query, bookingDetailID, bookingID).Scan(
        &passenger.BookingDetailID, &firstName, &lastName, &passenger.FlightClassName,
    )
    
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, fmt.Errorf("booking detail not found or doesn't belong to this booking")
        }
        return nil, fmt.Errorf("error validating ownership: %w", err)
    }
    
    passenger.PassengerName = fmt.Sprintf("%s %s", firstName, lastName)
    return &passenger, nil
}

// Helper function
func (r *checkinRepository) generateBoardingPassCode() string {
    return fmt.Sprintf("BP%d%d", time.Now().Unix()%10000, time.Now().Nanosecond()%1000)
}