package repository
import (
    "context"
    "errors"
    "log"
    "time" 
    "database/sql"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/dto"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
	"fmt"
)
type flightRepository struct {
	db *pgxpool.Pool
}
func NewFlightRepository(db *pgxpool.Pool) *flightRepository {
    return &flightRepository{db: db}
}
func (r *flightRepository) CreateFlight(flight *dto.Flight) (*dto.Flight, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // Bắt đầu transaction
    tx, err := r.db.Begin(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback(ctx)
    
    // Query insert flight
    query := `
        INSERT INTO flights (airline_id, aircraft_id, flight_number, departure_airport, arrival_airport, 
                           departure_time, arrival_time, duration_minutes, stops_count, tax_and_fees, 
                           status, gate, terminal, distance, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
        RETURNING flight_id
    `
    
    now := time.Now()
    
    var flightID int
    err = tx.QueryRow(ctx, query,
        flight.AirlineID, flight.AircraftID, flight.FlightNumber, flight.DepartureAirport,
        flight.ArrivalAirport, flight.DepartureTime, flight.ArrivalTime, flight.DurationMinutes,
        flight.StopsCount, flight.TaxAndFees, flight.Status, flight.Gate, flight.Terminal,
        flight.Distance, now, now).Scan(&flightID)
        
    if err != nil {
        return nil, fmt.Errorf("failed to create flight: %w", err)
    }
    
    flight.FlightID = flightID
    flight.CreatedAt = now
    flight.UpdatedAt = now
    
    // Commit transaction
    if err := tx.Commit(ctx); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    return flight, nil
}
func (r *flightRepository) CreateFlightClasses(flightID int, classes []*dto.FlightClass) ([]*dto.FlightClass, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    tx, err := r.db.Begin(ctx)
    if err != nil {
        log.Printf("Error starting transaction: %v", err)
        return nil, err
    }
    defer tx.Rollback(ctx)
    createdClasses := make([]*dto.FlightClass, 0, len(classes))
    for _, fc := range classes {
        query :=  `INSERT INTO flight_classes (flight_id, class, base_price, available_seats, total_seats)
            VALUES ($1, $2, $3, $4, $5)
            RETURNING flight_class_id
        `
        var flightClassID int
        err := tx.QueryRow(ctx, query,
             flightID, fc.Class, fc.BasePrice, fc.AvailableSeats, fc.TotalSeats).Scan(&flightClassID)
             
        if err != nil {
            log.Printf("Error creating flight class: %v", err)
            return nil, err
        }
        
        fc.FlightClassID = flightClassID
        fc.FlightID = flightID
        createdClasses = append(createdClasses, fc)
    }
    if err := tx.Commit(ctx); err != nil {
        log.Printf("Error committing transaction: %v", err)
        return nil, err
    }
    
    return createdClasses, nil
}

func (r *flightRepository) GetAll(page, limit int) ([]*dto.Flight, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    offset := (page - 1) * limit

    query := `
        SELECT flight_id, airline_id, aircraft_id, flight_number, departure_airport, arrival_airport,
               departure_time, arrival_time, duration_minutes, stops_count, tax_and_fees,
               total_seats, status, gate, terminal, distance, created_at, updated_at
        FROM flights
        ORDER BY departure_time
        LIMIT $1 OFFSET $2
    `

    rows, err := r.db.Query(ctx, query, limit, offset)
    if err != nil {
        log.Printf("Error querying flights: %v", err)
        return nil, err
    }
    defer rows.Close()

    flights := []*dto.Flight{}

    for rows.Next() {
        var flight dto.Flight
        err := rows.Scan(
            &flight.FlightID,
            &flight.AirlineID,
            &flight.AircraftID,
            &flight.FlightNumber,
            &flight.DepartureAirport,
            &flight.ArrivalAirport,
            &flight.DepartureTime,
            &flight.ArrivalTime,
            &flight.DurationMinutes,
            &flight.StopsCount,
            &flight.TaxAndFees,
            &flight.TotalSeats,
            &flight.Status,
            &flight.Gate,
            &flight.Terminal,
            &flight.Distance,
            &flight.CreatedAt,
            &flight.UpdatedAt,
        )
        if err != nil {
            log.Printf("Error scanning flight row: %v", err)
            return nil, err
        }
        flights = append(flights, &flight)
    }

    if err := rows.Err(); err != nil {
        log.Printf("Error iterating flight rows: %v", err)
        return nil, err
    }

    return flights, nil
}
func (r *flightRepository) Update(id int, flight *dto.Flight) (*dto.Flight, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `
        UPDATE flights
        SET airline_id = $1, aircraft_id = $2, flight_number = $3, departure_airport = $4, 
            arrival_airport = $5, departure_time = $6, arrival_time = $7, duration_minutes = $8,
            stops_count = $9, tax_and_fees = $10, total_seats = $11, status = $12,
            gate = $13, terminal = $14, distance = $15, updated_at = NOW()
        WHERE flight_id = $16
        RETURNING flight_id, airline_id, aircraft_id, flight_number, departure_airport, arrival_airport,
                  departure_time, arrival_time, duration_minutes, stops_count, tax_and_fees,
                  total_seats, status, gate, terminal, distance, created_at, updated_at
    `

    var updatedFlight dto.Flight
    err := r.db.QueryRow(ctx, query,
        flight.AirlineID,
        flight.AircraftID,
        flight.FlightNumber,
        flight.DepartureAirport,
        flight.ArrivalAirport,
        flight.DepartureTime,
        flight.ArrivalTime,
        flight.DurationMinutes,
        flight.StopsCount,
        flight.TaxAndFees,
        flight.TotalSeats,
        flight.Status,
        flight.Gate,
        flight.Terminal,
        flight.Distance,
        id,
    ).Scan(
        &updatedFlight.FlightID,
        &updatedFlight.AirlineID,
        &updatedFlight.AircraftID,
        &updatedFlight.FlightNumber,
        &updatedFlight.DepartureAirport,
        &updatedFlight.ArrivalAirport,
        &updatedFlight.DepartureTime,
        &updatedFlight.ArrivalTime,
        &updatedFlight.DurationMinutes,
        &updatedFlight.StopsCount,
        &updatedFlight.TaxAndFees,
        &updatedFlight.TotalSeats,
        &updatedFlight.Status,
        &updatedFlight.Gate,
        &updatedFlight.Terminal,
        &updatedFlight.Distance,
        &updatedFlight.CreatedAt,
        &updatedFlight.UpdatedAt,
    )

    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, errors.New("flight not found")
        }
        log.Printf("Error updating flight: %v", err)
        return nil, err
    }

    return &updatedFlight, nil
}
func (r *flightRepository) GetByFlightNumber(flightNumber string) (*dto.Flight, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `
        SELECT flight_id, airline_id, aircraft_id, flight_number, departure_airport, arrival_airport,
               departure_time, arrival_time, duration_minutes, stops_count, tax_and_fees,
               total_seats, status, gate, terminal, distance, created_at, updated_at
        FROM flights
        WHERE flight_number = $1
    `

    var flight dto.Flight
    err := r.db.QueryRow(ctx, query, flightNumber).Scan(
        &flight.FlightID,
        &flight.AirlineID,
        &flight.AircraftID,
        &flight.FlightNumber,
        &flight.DepartureAirport,
        &flight.ArrivalAirport,
        &flight.DepartureTime,
        &flight.ArrivalTime,
        &flight.DurationMinutes,
        &flight.StopsCount,
        &flight.TaxAndFees,
        &flight.TotalSeats,
        &flight.Status,
        &flight.Gate,
        &flight.Terminal,
        &flight.Distance,
        &flight.CreatedAt,
        &flight.UpdatedAt,
    )

    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, nil // Flight not found
        }
        log.Printf("Error getting flight by flight number: %v", err)
        return nil, err
    }

    return &flight, nil
}
func (r *flightRepository) GetByRoute(departureAirport, arrivalAirport string, date time.Time) ([]*dto.Flight, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Extract just the date part for comparison
    startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
    endDate := startDate.Add(24 * time.Hour)

    query := `
        SELECT flight_id, airline_id, aircraft_id, flight_number, departure_airport, arrival_airport,
               departure_time, arrival_time, duration_minutes, stops_count, tax_and_fees,
               total_seats, status, gate, terminal, distance, created_at, updated_at
        FROM flights
        WHERE departure_airport = $1 AND arrival_airport = $2
          AND departure_time >= $3 AND departure_time < $4
        ORDER BY departure_time
    `

    rows, err := r.db.Query(ctx, query, departureAirport, arrivalAirport, startDate, endDate)
    if err != nil {
        log.Printf("Error querying flights by route: %v", err)
        return nil, err
    }
    defer rows.Close()

    flights := []*dto.Flight{}

    for rows.Next() {
        var flight dto.Flight
        err := rows.Scan(
            &flight.FlightID,
            &flight.AirlineID,
            &flight.AircraftID,
            &flight.FlightNumber,
            &flight.DepartureAirport,
            &flight.ArrivalAirport,
            &flight.DepartureTime,
            &flight.ArrivalTime,
            &flight.DurationMinutes,
            &flight.StopsCount,
            &flight.TaxAndFees,
            &flight.TotalSeats,
            &flight.Status,
            &flight.Gate,
            &flight.Terminal,
            &flight.Distance,
            &flight.CreatedAt,
            &flight.UpdatedAt,
        )
        if err != nil {
            log.Printf("Error scanning flight row: %v", err)
            return nil, err
        }
        flights = append(flights, &flight)
    }

    if err := rows.Err(); err != nil {
        log.Printf("Error iterating flight rows: %v", err)
        return nil, err
    }

    return flights, nil
}
func (r *flightRepository) GetByID(id int) (*dto.Flight, []*dto.FlightClass, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
      flightQuery := `
        SELECT flight_id, airline_id, aircraft_id, flight_number, departure_airport, arrival_airport,
               departure_time, arrival_time, duration_minutes, stops_count, tax_and_fees,
               status, gate, terminal, distance, created_at, updated_at
        FROM flights
        WHERE flight_id = $1
    `
    var flight dto.Flight
    err := r.db.QueryRow(ctx, flightQuery, id).Scan(
        &flight.FlightID,
        &flight.AirlineID,
        &flight.AircraftID,
        &flight.FlightNumber,
        &flight.DepartureAirport,
        &flight.ArrivalAirport,
        &flight.DepartureTime,
        &flight.ArrivalTime,
        &flight.DurationMinutes,
        &flight.StopsCount,
        &flight.TaxAndFees,
        &flight.Status,
        &flight.Gate,
        &flight.Terminal,
        &flight.Distance,
        &flight.CreatedAt,
        &flight.UpdatedAt,
    )
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, nil, nil 
        }
        return nil, nil, fmt.Errorf("error getting flight by ID: %w", err)
    }
   flightClassesQuery := `
        SELECT flight_class_id, flight_id, class, base_price, available_seats, total_seats
        FROM flight_classes
        WHERE flight_id = $1
        ORDER BY base_price ASC
    `

    rows, err := r.db.Query(ctx, flightClassesQuery, id)
    if err != nil {
        return &flight, nil, fmt.Errorf("error querying flight classes: %w", err)
    }
    defer rows.Close()

    flightClasses := []*dto.FlightClass{}

    for rows.Next() {
        var fc dto.FlightClass
        if err := rows.Scan(
            &fc.FlightClassID,
            &fc.FlightID,
            &fc.Class,
            &fc.BasePrice,
            &fc.AvailableSeats,
            &fc.TotalSeats,
        ); err != nil {
            return &flight, nil, fmt.Errorf("error scanning flight class: %w", err)
        }
        flightClasses = append(flightClasses, &fc)
    }

    if err := rows.Err(); err != nil {
        return &flight, nil, fmt.Errorf("error iterating flight class rows: %w", err)
    }

    return &flight, flightClasses, nil
}
func (r *flightRepository) GetByAirline(airlineID int, page, limit int) ([]*dto.Flight, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    offset := (page - 1) * limit

    query := `
        SELECT flight_id, airline_id, aircraft_id, flight_number, departure_airport, arrival_airport,
               departure_time, arrival_time, duration_minutes, stops_count, tax_and_fees,
               total_seats, status, gate, terminal, distance, created_at, updated_at
        FROM flights
        WHERE airline_id = $1
        ORDER BY departure_time
        LIMIT $2 OFFSET $3
    `

    rows, err := r.db.Query(ctx, query, airlineID, limit, offset)
    if err != nil {
        log.Printf("Error querying flights by airline: %v", err)
        return nil, err
    }
    defer rows.Close()

    flights := []*dto.Flight{}

    for rows.Next() {
        var flight dto.Flight
        err := rows.Scan(
            &flight.FlightID,
            &flight.AirlineID,
            &flight.AircraftID,
            &flight.FlightNumber,
            &flight.DepartureAirport,
            &flight.ArrivalAirport,
            &flight.DepartureTime,
            &flight.ArrivalTime,
            &flight.DurationMinutes,
            &flight.StopsCount,
            &flight.TaxAndFees,
            &flight.TotalSeats,
            &flight.Status,
            &flight.Gate,
            &flight.Terminal,
            &flight.Distance,
            &flight.CreatedAt,
            &flight.UpdatedAt,
        )
        if err != nil {
            log.Printf("Error scanning flight row: %v", err)
            return nil, err
        }
        flights = append(flights, &flight)
    }

    if err := rows.Err(); err != nil {
        log.Printf("Error iterating flight rows: %v", err)
        return nil, err
    }

    return flights, nil
}

func (r *flightRepository) GetByStatus(status string, page, limit int) ([]*dto.Flight, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    offset := (page - 1) * limit

    query := `
        SELECT flight_id, airline_id, aircraft_id, flight_number, departure_airport, arrival_airport,
               departure_time, arrival_time, duration_minutes, stops_count, tax_and_fees,
               total_seats, status, gate, terminal, distance, created_at, updated_at
        FROM flights
        WHERE status = $1
        ORDER BY departure_time
        LIMIT $2 OFFSET $3
    `

    rows, err := r.db.Query(ctx, query, status, limit, offset)
    if err != nil {
        log.Printf("Error querying flights by status: %v", err)
        return nil, err
    }
    defer rows.Close()

    flights := []*dto.Flight{}

    for rows.Next() {
        var flight dto.Flight
        err := rows.Scan(
            &flight.FlightID,
            &flight.AirlineID,
            &flight.AircraftID,
            &flight.FlightNumber,
            &flight.DepartureAirport,
            &flight.ArrivalAirport,
            &flight.DepartureTime,
            &flight.ArrivalTime,
            &flight.DurationMinutes,
            &flight.StopsCount,
            &flight.TaxAndFees,
            &flight.TotalSeats,
            &flight.Status,
            &flight.Gate,
            &flight.Terminal,
            &flight.Distance,
            &flight.CreatedAt,
            &flight.UpdatedAt,
        )
        if err != nil {
            log.Printf("Error scanning flight row: %v", err)
            return nil, err
        }
        flights = append(flights, &flight)
    }

    if err := rows.Err(); err != nil {
        log.Printf("Error iterating flight rows: %v", err)
        return nil, err
    }

    return flights, nil
}

func (r *flightRepository) SearchFlights(
    departureAirport string,
    arrivalAirport string, 
    departureDate time.Time, 
    class string, 
    airlineIDs []int, 
    maxStops int, 
    page int, 
    limit int,
    sortBy string,
    sortOrder string,
) ([]*dto.FlightSearchResult, error) {
    // Thêm log để debug
    log.Printf("Tìm kiếm chuyến bay: từ %s đến %s vào ngày %s, hạng ghế: %s", 
        departureAirport, arrivalAirport, departureDate.Format("2006-01-02"), class)
    
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // Tạo câu truy vấn cơ bản với thêm log
    baseQuery := `
        SELECT f.flight_id, f.airline_id, a.name as airline_name, f.flight_number, 
               f.departure_airport, f.arrival_airport, f.departure_time, f.arrival_time, 
               f.duration_minutes, f.stops_count, f.tax_and_fees, 
               f.status, f.gate, f.terminal, f.distance, 
               fc.class as flight_class, fc.base_price as class_price,
               fc.available_seats as class_availability, f.total_seats
        FROM flights f
        JOIN flight_classes fc ON f.flight_id = fc.flight_id
        JOIN airlines a ON f.airline_id = a.airline_id
        WHERE f.departure_airport = $1 
        AND f.arrival_airport = $2 
        AND f.departure_time::date = $3::date
        AND fc.class = $4
        AND fc.available_seats > 0
    `
    
    // Log câu truy vấn và tham số
    log.Printf("Câu truy vấn SQL: %s", baseQuery)
    log.Printf("Tham số: %s, %s, %s, %s", 
        departureAirport, arrivalAirport, departureDate.Format("2006-01-02"), class)
    
    // Thực thi truy vấn và kiểm tra kết quả
    rows, err := r.db.Query(ctx, baseQuery, departureAirport, arrivalAirport, departureDate, class)
    if err != nil {
        log.Printf("Lỗi khi thực thi truy vấn tìm kiếm chuyến bay: %v", err)
        return nil, fmt.Errorf("lỗi truy vấn chuyến bay: %w", err)
    }
    defer rows.Close()
    
    // Xử lý kết quả với log bổ sung
    results := []*dto.FlightSearchResult{}
    rowCount := 0
    
    for rows.Next() {
        rowCount++
        var result dto.FlightSearchResult
        var totalSeats sql.NullInt32
        
        err := rows.Scan(
            &result.FlightID,
            &result.AirlineID,
            &result.AirlineName,
            &result.FlightNumber,
            &result.DepartureAirport,
            &result.ArrivalAirport,
            &result.DepartureTime,
            &result.ArrivalTime,
            &result.DurationMinutes,
            &result.StopsCount,
            &result.TaxAndFees,
            &result.Status,
            &result.Gate,
            &result.Terminal,
            &result.Distance,
            &result.FlightClass,
            &result.ClassPrice,
            &result.ClassAvailability,
            &totalSeats,
        )
        
        if err != nil {
            log.Printf("Lỗi khi quét dữ liệu chuyến bay: %v", err)
            continue  // Bỏ qua hàng lỗi và tiếp tục
        }
        
        // Xử lý giá trị NULL
        if totalSeats.Valid {
            result.TotalSeats = int(totalSeats.Int32)
        }
        
        // Tính tổng giá
        result.TotalPrice = result.ClassPrice + result.TaxAndFees
        
        results = append(results, &result)
    }
    
    log.Printf("Tìm thấy %d chuyến bay phù hợp với tiêu chí", rowCount)
    
    if err := rows.Err(); err != nil {
        log.Printf("Lỗi khi lặp qua kết quả tìm kiếm: %v", err)
        return nil, fmt.Errorf("lỗi khi lặp qua chuyến bay: %w", err)
    }
    
    return results, nil
}
func (r *flightRepository) Count() (int, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var count int
    query := `SELECT COUNT(*) FROM flights`

    err := r.db.QueryRow(ctx, query).Scan(&count)
    if err != nil {
        log.Printf("Error counting flights: %v", err)
        return 0, err
    }

    return count, nil
}

func (r *flightRepository) CountBySearch(departureAirport, arrivalAirport string, departureDate time.Time) (int, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Extract just the date part for comparison
    startDate := time.Date(departureDate.Year(), departureDate.Month(), departureDate.Day(), 0, 0, 0, 0, departureDate.Location())
    endDate := startDate.Add(24 * time.Hour)

    var count int
    query := `
        SELECT COUNT(*) 
        FROM flights
        WHERE departure_airport = $1 AND arrival_airport = $2
          AND departure_time >= $3 AND departure_time < $4
    `

    err := r.db.QueryRow(ctx, query, departureAirport, arrivalAirport, startDate, endDate).Scan(&count)
    if err != nil {
        log.Printf("Error counting flights by search: %v", err)
        return 0, err
    }

    return count, nil
}

func (r *flightRepository) UpdateStatus(id int, status string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `UPDATE flights SET status = $1, updated_at = NOW() WHERE flight_id = $2`
    
    commandTag, err := r.db.Exec(ctx, query, status, id)
    if err != nil {
        log.Printf("Error updating flight status: %v", err)
        return err
    }

    if commandTag.RowsAffected() == 0 {
        return errors.New("flight not found")
    }

    return nil
}