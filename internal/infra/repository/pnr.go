package repository

import (
    "context"
    "errors"
    "fmt"
    "math/rand"
    "time"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/dto"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type pnrRepository struct {
    db *pgxpool.Pool
}

func NewPNRRepository(db *pgxpool.Pool) *pnrRepository {
    return &pnrRepository{db: db}
}
func (r *pnrRepository) GeneratePNR(bookingID int64) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    for attempts := 0; attempts < 10; attempts++ {
        pnrCode := r.generateRandomPNRCode()

        var exists bool
        err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM pnr WHERE pnr_code = $1)", pnrCode).Scan(&exists)
        if err != nil {
            return "", fmt.Errorf("failed to check PNR existence: %v", err)
        }

        if !exists {
            return pnrCode, nil
        }
    }

    return "", errors.New("failed to generate unique PNR after multiple attempts")
}
func (r *pnrRepository) generateRandomPNRCode() string {

    const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    prefix := make([]byte, 2)
    for i := range prefix {
        prefix[i] = letters[rand.Intn(len(letters))]
    }

    digits := rand.Intn(10000)
    
    return fmt.Sprintf("%s%04d", string(prefix), digits)
}
func (r *pnrRepository) CreatePnr(pnr *dto.PNR) (*dto.PNR, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if pnr.PNRCode == "" {
        var bookingID int64
        _, err := fmt.Sscanf(pnr.BookingID, "%d", &bookingID)
        if err != nil {
            return nil, fmt.Errorf("invalid booking ID format: %v", err)
        }

        // Sử dụng hàm GeneratePNR có sẵn
        generatedCode, err := r.GeneratePNR(bookingID)
        if err != nil {
            return nil, fmt.Errorf("failed to generate PNR code: %v", err)
        }
        pnr.PNRCode = generatedCode
    }

    currentTime := time.Now().Format("2006-01-02 15:04:05")
    if pnr.IssuedAt == "" {
        pnr.IssuedAt = currentTime
    }
    if pnr.Status == "" {
        pnr.Status = "active"
    }
    if pnr.ExpiresAt == "" {
  
        expiryTime := time.Now().AddDate(1, 0, 0).Format("2006-01-02 15:04:05")
        pnr.ExpiresAt = expiryTime
    }

    query := `
        INSERT INTO pnrs (
            pnr_code, booking_id, flight_id, return_flight_id,
            status, issued_at, expires_at, created_by, pnr_data
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING pnr_code, booking_id, flight_id, return_flight_id,
                  status, issued_at, expires_at, created_by, modified_by, pnr_data
    `

    var result dto.PNR
    err := r.db.QueryRow(ctx, query,
        pnr.PNRCode,
        pnr.BookingID,
        pnr.FlightID,
        pnr.ReturnFlightID,
        pnr.Status,
        pnr.IssuedAt,
        pnr.ExpiresAt,
        pnr.CreatedBy,
        pnr.PNRData,
    ).Scan(
        &result.PNRCode,
        &result.BookingID,
        &result.FlightID,
        &result.ReturnFlightID,
        &result.Status,
        &result.IssuedAt,
        &result.ExpiresAt,
        &result.CreatedBy,
        &result.ModifiedBy,
        &result.PNRData,
    )

    if err != nil {
        return nil, fmt.Errorf("failed to create PNR: %v", err)
    }

    return &result, nil
}
func (r *pnrRepository) CheckPnrExists(code string) (bool, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `SELECT EXISTS(SELECT 1 FROM pnrs WHERE pnr_code = $1)`

    var exists bool
    err := r.db.QueryRow(ctx, query, code).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("failed to check PNR existence: %v", err)
    }

    return exists, nil
}
func (r *pnrRepository) GetPnrByBookingID(bookingID int64) (*dto.PNR, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    bookingIDStr := fmt.Sprintf("%d", bookingID)

    query := `
        SELECT pnr_code, booking_id, flight_id, return_flight_id,
               status, issued_at, expires_at, created_by, modified_by, pnr_data
        FROM pnrs
        WHERE booking_id = $1
    `

    var pnr dto.PNR
    err := r.db.QueryRow(ctx, query, bookingIDStr).Scan(
        &pnr.PNRCode,
        &pnr.BookingID,
        &pnr.FlightID,
        &pnr.ReturnFlightID,
        &pnr.Status,
        &pnr.IssuedAt,
        &pnr.ExpiresAt,
        &pnr.CreatedBy,
        &pnr.ModifiedBy,
        &pnr.PNRData,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, fmt.Errorf("no PNR found for booking ID %d", bookingID)
        }
        return nil, fmt.Errorf("failed to get PNR: %v", err)
    }

    return &pnr, nil
}

func (r *pnrRepository) GetBookingIDByPnrCode(pnrCode string) (int64, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `SELECT booking_id FROM pnrs WHERE pnr_code = $1`

    var bookingIDStr string
    err := r.db.QueryRow(ctx, query, pnrCode).Scan(&bookingIDStr)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return 0, fmt.Errorf("no booking found with PNR code %s", pnrCode)
        }
        return 0, fmt.Errorf("failed to get booking ID: %v", err)
    }

    var bookingID int64
    _, err = fmt.Sscanf(bookingIDStr, "%d", &bookingID)
    if err != nil {
        return 0, fmt.Errorf("invalid booking ID format: %v", err)
    }

    return bookingID, nil
}

func (r *pnrRepository) UpdatePnr(pnr *dto.PNR) (*dto.PNR, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `
        UPDATE pnrs
        SET status = $1, 
            expires_at = $2, 
            modified_by = $3,
            pnr_data = $4,
            flight_id = $5,
            return_flight_id = $6
        WHERE pnr_code = $7
        RETURNING pnr_code, booking_id, flight_id, return_flight_id,
                  status, issued_at, expires_at, created_by, modified_by, pnr_data
    `

    var result dto.PNR
    err := r.db.QueryRow(ctx, query,
        pnr.Status,
        pnr.ExpiresAt,
        pnr.ModifiedBy,
        pnr.PNRData,
        pnr.FlightID,
        pnr.ReturnFlightID,
        pnr.PNRCode,
    ).Scan(
        &result.PNRCode,
        &result.BookingID,
        &result.FlightID,
        &result.ReturnFlightID,
        &result.Status,
        &result.IssuedAt,
        &result.ExpiresAt,
        &result.CreatedBy,
        &result.ModifiedBy,
        &result.PNRData,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, fmt.Errorf("PNR with code %s not found", pnr.PNRCode)
        }
        return nil, fmt.Errorf("failed to update PNR: %v", err)
    }

    return &result, nil
}