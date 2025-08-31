package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Khangvn20/FlyJourney_Backend/internal/core/common/utils"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/dto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
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
    
    tx, err := r.db.Begin(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback(ctx)
    
     query := `
        INSERT INTO flights (airline_id, flight_number, departure_airport, arrival_airport,
                           duration_minutes, stops_count, tax_and_fees, total_seats, status,
                           distance, departure_time, arrival_time, 
                           departure_airport_code, arrival_airport_code, currency,
                           created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
        RETURNING flight_id
    `
    
    now := time.Now()
    
    var flightID int
    err = tx.QueryRow(ctx, query,
        flight.AirlineID,          
        flight.FlightNumber,         
        flight.DepartureAirport,     
        flight.ArrivalAirport,      
        flight.DurationMinutes,      
        flight.StopsCount,           
        flight.TaxAndFees,        
        flight.TotalSeats,         
        flight.Status,             
        flight.Distance,           
        flight.DepartureTime,       
        flight.ArrivalTime,         
        flight.DepartureAirportCode, 
        flight.ArrivalAiportCode,   
        flight.Currency,             
        now,
        now).Scan(&flightID)
        
    if err != nil {
        return nil, fmt.Errorf("failed to create flight: %w", err)
    }
    
    flight.FlightID = flightID
    flight.CreatedAt = now
    flight.UpdatedAt = now
    
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
    totalSeats := 0
    for _, fc := range classes {
        query := `
            INSERT INTO flight_classes (flight_id, class, base_price, available_seats, total_seats, base_price_child, base_price_infant,fare_class_code)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
            RETURNING flight_class_id
        `
        
        var flightClassID int
        err := tx.QueryRow(ctx, query,
            flightID, fc.Class, fc.BasePrice, fc.AvailableSeats, fc.TotalSeats,fc.BasePriceChild,fc.BasePriceInfant,fc.FareClassCode).Scan(&flightClassID)
             
        if err != nil {
            log.Printf("Error creating flight class: %v", err)
            return nil, err
        }
        
        fc.FlightClassID = flightClassID
        fc.FlightID = flightID
        createdClasses = append(createdClasses, fc)
        totalSeats += fc.TotalSeats
    }
    updateQuery := `
        UPDATE flights 
        SET total_seats = $1 
        WHERE flight_id = $2
    `
    
    _, err = tx.Exec(ctx, updateQuery, totalSeats, flightID)
    if err != nil {
        log.Printf("Error updating flight total_seats: %v", err)
        return nil, err
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
        SELECT flight_id, airline_id, flight_number, departure_airport, arrival_airport,
               departure_time, arrival_time, duration_minutes, stops_count, tax_and_fees,
               total_seats, status, distance, departure_airport_code, arrival_airport_code,
               currency, created_at, updated_at
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
            &flight.Distance,
            &flight.DepartureAirportCode,
            &flight.ArrivalAiportCode,
            &flight.Currency,
            &flight.CreatedAt,
            &flight.UpdatedAt,
        )
        if err != nil {
            log.Printf("Error scanning flight row: %v", err)
            return nil, err
        }

        flights = append(flights, &flight)
    }

    return flights, nil
}
func (r *flightRepository) Update(id int, flight *dto.Flight) (*dto.Flight, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `
        UPDATE flights
        SET airline_id = $1, flight_number = $2, departure_airport = $3, 
            arrival_airport = $4, departure_time = $5, arrival_time = $6, 
            duration_minutes = $7, stops_count = $8, tax_and_fees = $9, 
            total_seats = $10, status = $11, distance = $12,
            departure_airport_code = $13, arrival_airport_code = $14, 
            currency = $15, updated_at = NOW()
        WHERE flight_id = $16
        RETURNING flight_id, airline_id, flight_number, departure_airport, arrival_airport,
                  departure_time, arrival_time, duration_minutes, stops_count, tax_and_fees,
                  total_seats, status, distance, departure_airport_code, arrival_airport_code,
                  currency, created_at, updated_at
    `

    var updatedFlight dto.Flight
    err := r.db.QueryRow(ctx, query,
        flight.AirlineID,         
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
        flight.Distance,              
        flight.DepartureAirportCode,   
        flight.ArrivalAiportCode,      
        flight.Currency,               
        id,                            
    ).Scan(
        &updatedFlight.FlightID,
        &updatedFlight.AirlineID,
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
        &updatedFlight.Distance,
        &updatedFlight.DepartureAirportCode,
        &updatedFlight.ArrivalAiportCode,
        &updatedFlight.Currency,
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
        SELECT flight_id, airline_id, flight_number, departure_airport, arrival_airport,
               departure_time, arrival_time, duration_minutes, stops_count, tax_and_fees,
               total_seats, status, distance, created_at, updated_at
        FROM flights
        WHERE flight_number = $1
    `

    var flight dto.Flight
    err := r.db.QueryRow(ctx, query, flightNumber).Scan(
        &flight.FlightID,
        &flight.AirlineID,
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
        &flight.Distance,
        &flight.CreatedAt,
        &flight.UpdatedAt,
    )

    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, nil 
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
        SELECT flight_id, airline_id, flight_number, departure_airport, arrival_airport,
               departure_time, arrival_time, duration_minutes, stops_count, tax_and_fees,
               total_seats, status, distance, departure_airport_code, arrival_airport_code,
               currency, created_at, updated_at
        FROM flights
        WHERE flight_id = $1
    `
    
    var flight dto.Flight
    err := r.db.QueryRow(ctx, flightQuery, id).Scan(
        &flight.FlightID,
        &flight.AirlineID,
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
        &flight.Distance,
        &flight.DepartureAirportCode,
        &flight.ArrivalAiportCode,
        &flight.Currency,
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
        SELECT fc.flight_class_id, fc.flight_id, fc.class, fc.base_price, fc.available_seats, 
               fc.total_seats, fc.base_price_child, fc.base_price_infant, fc.fare_class_code,
               fcc.fare_class_code,fcc.cabin_class, fcc.refundable, fcc.changeable, fcc.baggage_kg, fcc.refund_change_policy,fcc.description
        FROM flight_classes fc
        LEFT JOIN fare_classes fcc ON fc.fare_class_code = fcc.fare_class_code
        WHERE fc.flight_id = $1
        ORDER BY fc.base_price ASC
    `

    rows, err := r.db.Query(ctx, flightClassesQuery, id)
    if err != nil {
        return &flight, nil, fmt.Errorf("error querying flight classes: %w", err)
    }
    defer rows.Close()

    flightClasses := []*dto.FlightClass{}

    for rows.Next() {
        var fc dto.FlightClass
        var fareClass dto.FareClasses
        
        if err := rows.Scan(
            &fc.FlightClassID,
            &fc.FlightID,
            &fc.Class,
            &fc.BasePrice,
            &fc.AvailableSeats,
            &fc.TotalSeats,
            &fc.BasePriceChild,
            &fc.BasePriceInfant,
            &fc.FareClassCode,
            &fareClass.FareClassCode,
            &fareClass.CabinClass,
            &fareClass.Refundable,
            &fareClass.Changeable,
            &fareClass.Baggage_kg,
            &fareClass.RefundChangePolicy,
            &fareClass.Description,
        ); err != nil {
            return &flight, nil, fmt.Errorf("error scanning flight class: %w", err)
        }
         fc.FareClassDetails = &fareClass
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
    departureDate string,
    flightClass string,
    airlineIDs []int,
    maxStops int,
    page int,
    limit int,
    sortBy string,
    sortOrder string,
    forUser bool,
) ([]*dto.FlightSearchResult, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    offset := (page - 1) * limit
    parsedDate, err := utils.ParseTime(departureDate)
    if err != nil {
        log.Printf("Error parsing departure date: %v", err)
        return nil, fmt.Errorf("invalid departure date format. Required format: dd/mm/yyyy: %w", err)
    }
    
    // Format date for PostgreSQL DATE comparison (YYYY-MM-DD)
    formattedDate := parsedDate.Format("2006-01-02")
    
    log.Printf("Original date: %s, Parsed date: %s, Formatted for DB: %s", 
        departureDate, parsedDate.String(), formattedDate)



    // Base query
    query := `
        SELECT 
            f.flight_id,
            f.flight_number,
            f.airline_id,
            a.name as airline_name,
            a.logo_url,
            f.departure_airport_code,
            f.arrival_airport_code,
            f.departure_airport,
            f.arrival_airport,
            f.departure_time,
            f.arrival_time,
            f.duration_minutes,
            f.stops_count,
            f.distance,
            fc.class as flight_class,
            f.total_seats,
            fc.flight_class_id,
            fc.base_price,
            fc.base_price_child,
            fc.base_price_infant,
            f.tax_and_fees,
            f.currency,
            fcc.fare_class_code,
            fcc.cabin_class,
            fcc.refundable,
            fcc.changeable,
            fcc.refund_change_policy,
            fcc.baggage_kg,
            fcc.description
        FROM flights f
        INNER JOIN flight_classes fc ON f.flight_id = fc.flight_id
        LEFT JOIN airlines a ON f.airline_id = a.airline_id
        LEFT JOIN fare_classes fcc ON fc.fare_class_code = fcc.fare_class_code
        WHERE f.departure_airport_code = $1
        AND f.arrival_airport_code = $2
        AND DATE(f.departure_time) = $3`


    args := []interface{}{departureAirport, arrivalAirport, formattedDate}
    argIndex := 4

    
    // Add airline filter
    if len(airlineIDs) > 0 {
        query += fmt.Sprintf(" AND f.airline_id = ANY($%d)", argIndex)
        args = append(args, pq.Array(airlineIDs))
        argIndex++
    }

    // Add stops filter
    if maxStops >= 0 {
        query += fmt.Sprintf(" AND f.stops_count <= $%d", argIndex)
        args = append(args, maxStops)
        argIndex++
    }

    // Add status filter for users
    if forUser {
        query += " AND f.status IN ('scheduled', 'boarding')"
        query += " AND fc.available_seats > 0"
    }

    // Add sorting
    validSortFields := map[string]string{
        "departure_time": "f.departure_time",
        "arrival_time":   "f.arrival_time",
        "price":          "fc.base_price",
        "duration":       "f.duration_minutes",
        "stops":          "f.stops_count",
    }

    if sortField, ok := validSortFields[sortBy]; ok {
        if sortOrder == "DESC" || sortOrder == "desc" {
            query += fmt.Sprintf(" ORDER BY %s DESC", sortField)
        } else {
            query += fmt.Sprintf(" ORDER BY %s ASC", sortField)
        }
    } else {
        query += " ORDER BY f.departure_time ASC"
    }

    // Add pagination
    query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
    args = append(args, limit, offset)

    log.Printf("Search query: %s", query)
    log.Printf("Search args: %v", args)

    rows, err := r.db.Query(ctx, query, args...)
    if err != nil {
        log.Printf("Error searching flights: %v", err)
        return nil, err
    }
    defer rows.Close()

    results := []*dto.FlightSearchResult{}

    for rows.Next() {
        var result dto.FlightSearchResult
        var fareClass dto.FareClasses
        var airlineName sql.NullString
        var fareClassCode, cabinClass, baggageKg, description sql.NullString
        var refundable, changeable sql.NullBool

        err := rows.Scan(
            &result.FlightID,
            &result.FlightNumber,
            &result.AirlineID,
            &airlineName,
            &result.LogoUrl,
            &result.DepartureAirportCode,
            &result.ArrivalAirportCode,
            &result.DepartureAirport,
            &result.ArrivalAirport,
            &result.DepartureTime,
            &result.ArrivalTime,
            &result.Duration,
            &result.StopsCount,
            &result.Distance,
            &result.FlightClass,
            &result.TotalSeats,
            &result.FlightClassID,
            &result.Pricing.BasePrices.Adult,
            &result.Pricing.BasePrices.Child,
            &result.Pricing.BasePrices.Infant,
            &result.TaxAndFees,
            &result.Pricing.Currency,
            &fareClassCode,
            &cabinClass,
            &refundable,
            &changeable,
            &fareClass.RefundChangePolicy,
            &baggageKg,
            &description,
        )
        if err != nil {
            log.Printf("Error scanning flight search result: %v", err)
            return nil, err
        }

        // Set airline name
        if airlineName.Valid {
            result.AirlineName = airlineName.String
        }
        if fareClassCode.Valid {
            fareClass.FareClassCode = fareClassCode.String
            fareClass.CabinClass = cabinClass.String
            fareClass.Refundable = refundable.Bool
            fareClass.Changeable = changeable.Bool
            fareClass.Baggage_kg = baggageKg.String
            fareClass.Description = description.String
            result.FareClassDetails = &fareClass
        }

        // Set fare class details
        if fareClassCode.Valid {
            fareClass.FareClassCode = fareClassCode.String
            fareClass.CabinClass = cabinClass.String
            fareClass.Refundable = refundable.Bool
            fareClass.Changeable = changeable.Bool
            fareClass.Baggage_kg = baggageKg.String
            fareClass.Description = description.String
            result.FareClassDetails = &fareClass
        }

        // ✅ Calculate taxes and prices for adults (always required)
        adultTax := result.TaxAndFees
        result.Pricing.Taxes.Adult = adultTax
        result.Pricing.TotalPrices.Adult = result.Pricing.BasePrices.Adult + adultTax

        // ✅ Only calculate child pricing if child base price > 0
        if result.Pricing.BasePrices.Child > 0 {
            childTax := result.TaxAndFees * 0.75 // 75% of adult tax
            result.Pricing.Taxes.Child = childTax
            result.Pricing.TotalPrices.Child = result.Pricing.BasePrices.Child + childTax
        }

        // ✅ Only calculate infant pricing if infant base price > 0  
        if result.Pricing.BasePrices.Infant > 0 {
            infantTax := result.TaxAndFees * 0.25 // 25% of adult tax
            result.Pricing.Taxes.Infant = infantTax
            result.Pricing.TotalPrices.Infant = result.Pricing.BasePrices.Infant + infantTax
        }

        // Calculate grand total (for 1 adult by default - service sẽ override)
        result.Pricing.GrandTotal = result.Pricing.TotalPrices.Adult

        results = append(results, &result)
    }
    if err := rows.Err(); err != nil {
        log.Printf("Error iterating flight search results: %v", err)
        return nil, err
    }

    return results, nil
}
func (r *flightRepository) Count() (int, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `SELECT COUNT(*) FROM flights`
    
    var count int
    err := r.db.QueryRow(ctx, query).Scan(&count)
    if err != nil {
        log.Printf("Error counting flights: %v", err)
        return 0, err
    }

    return count, nil
}

func (r *flightRepository) CountBySearch(departureAirport, arrivalAirport string, departureDate string, forUser bool) (int, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Parse and format departure date for database query
    parsedDate, err := utils.ParseTime(departureDate)
    if err != nil {
        log.Printf("Error parsing departure date: %v", err)
        return 0, fmt.Errorf("invalid departure date format: %w", err)
    }
    
    // Format date for PostgreSQL DATE comparison
    formattedDate := parsedDate.Format("2006-01-02")

    query := `
        SELECT COUNT(*)
        FROM flights f
        JOIN flight_classes fc ON f.flight_id = fc.flight_id
        WHERE f.departure_airport_code = $1
        AND f.arrival_airport_code = $2
        AND DATE(f.departure_time) = $3`

    if forUser {
        query += " AND f.status IN ('scheduled', 'boarding')"
        query += " AND fc.available_seats > 0"
    }

    var count int
    err = r.db.QueryRow(ctx, query, departureAirport, arrivalAirport, formattedDate).Scan(&count)
    if err != nil {
        log.Printf("Error counting flights by search: %v", err)
        return 0, err
    }

    return count, nil
}
func (r *flightRepository) CountByDate (departureTime time.Time) (int, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    startDate := time.Date(departureTime.Year(), departureTime.Month(), departureTime.Day(), 0, 0, 0, 0, departureTime.Location())
    endDate := startDate.Add(24 * time.Hour)

    query := `
        SELECT COUNT(*)
        FROM flights
        WHERE departure_time >= $1 AND departure_time < $2
    `

    var count int
    err := r.db.QueryRow(ctx, query, startDate, endDate).Scan(&count)
    if err != nil {
        log.Printf("Error counting flights by date: %v", err)
        return 0, err
    }

    return count, nil
}
func (r *flightRepository) GetFareClassesByFlightID(flightID int) ([]*dto.FareClasses, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    query := `
        SELECT DISTINCT
            fcc.fare_class_code,
            fcc.cabin_class, 
            fcc.refundable, 
            fcc.changeable, 
            fcc.baggage_kg, 
            fcc.refund_change_policy,
            fcc.description
        FROM flight_classes fc
        INNER JOIN fare_classes fcc ON fc.fare_class_code = fcc.fare_class_code
        WHERE fc.flight_id = $1
        ORDER BY fcc.fare_class_code
    `

    rows, err := r.db.Query(ctx, query, flightID)
    if err != nil {
        log.Printf("Error querying fare classes detail by flight ID: %v", err)
        return nil, fmt.Errorf("error querying fare classes detail: %w", err)
    }
    defer rows.Close()

    fareClasses := []*dto.FareClasses{}

    for rows.Next() {
        var fareClass dto.FareClasses

        err := rows.Scan(
            &fareClass.FareClassCode,
            &fareClass.CabinClass,
            &fareClass.Refundable,
            &fareClass.Changeable,
            &fareClass.Baggage_kg,
            &fareClass.RefundChangePolicy,
            &fareClass.Description,
        )
        if err != nil {
            log.Printf("Error scanning fare class detail row: %v", err)
            return nil, fmt.Errorf("error scanning fare class detail: %w", err)
        }

        fareClasses = append(fareClasses, &fareClass)
    }

    if err := rows.Err(); err != nil {
        log.Printf("Error iterating fare class detail rows: %v", err)
        return nil, fmt.Errorf("error iterating detail rows: %w", err)
    }

    return fareClasses, nil
}
func (r *flightRepository) GetFlightsByDate(departureTime time.Time, page, limit int ) ([]*dto.Flight, error) {
      ctx, cancel :=context.WithTimeout(context.Background(), 5*time.Second)
      defer cancel()
      offset := (page - 1) * limit
      startDate := time.Date(departureTime.Year(), departureTime.Month(), departureTime.Day(), 0, 0, 0, 0, departureTime.Location())
      endDate := startDate.Add(24 * time.Hour)
      query := `
        SELECT DiSTINCT f.flight_id, f.airline_id, f.flight_number, f.departure_airport, f.arrival_airport,
        f.departure_time, f.arrival_time, f.duration_minutes, f.stops_count, f.tax_and_fees,f.total_seats,
        f.status, f.distance, f.departure_airport_code, f.arrival_airport_code, f.currency, f.created_at, f.updated_at,
        fc.flight_class_id, fc.class, fc.base_price, fc.available_seats, fc.total_seats,fc.base_price_child, fc.base_price_infant, fc.fare_class_code,
        fcc.fare_class_code, fcc.cabin_class, fcc.refundable, fcc.changeable, fcc.baggage_kg,fcc.refund_change_policy ,fcc.description
        FROM flights f
        INNER JOIN flight_classes fc ON f.flight_id = fc.flight_id
        LEFT JOIN fare_classes fcc ON fc.fare_class_code = fcc.fare_class_code
        WHERE f.departure_time >= $1 AND f.departure_time < $2
        ORDER BY f.departure_time
        LIMIT $3 OFFSET $4`
    rows, err := r.db.Query(ctx, query, startDate, endDate, limit, offset)
    if err != nil {
        log.Printf("Error querying flights by date: %v", err)
        return nil, err
    }
    defer rows.Close()
    flights := []*dto.Flight{}
    for rows.Next() {
        var flight dto.Flight
        var flightClass dto.FlightClass
        var fareClass dto.FareClasses

        err := rows.Scan(
            &flight.FlightID,
            &flight.AirlineID,
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
            &flight.Distance,
            &flight.DepartureAirportCode,
            &flight.ArrivalAiportCode,
            &flight.Currency,
            &flight.CreatedAt,
            &flight.UpdatedAt,
            &flightClass.FlightClassID,
            &flightClass.Class,
            &flightClass.BasePrice,
            &flightClass.AvailableSeats,
            &flightClass.TotalSeats,
            &flightClass.BasePriceChild,
            &flightClass.BasePriceInfant,
            &flightClass.FareClassCode,
            &fareClass.FareClassCode,
            &fareClass.CabinClass,
            &fareClass.Refundable,
            &fareClass.Changeable,
            &fareClass.Baggage_kg,
            &fareClass.RefundChangePolicy, 
            &fareClass.Description,
        )
        if err != nil {
            log.Printf("Error scanning flight row: %v", err)
            return nil, err
        }
        flight.FlightClasses = append(flight.FlightClasses, &flightClass)
        flightClass.FareClassDetails = &fareClass
        flights = append(flights, &flight)
    }
    return flights, nil
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
func (r *flightRepository) SearchRoundtripFlights(
    departureAirport string,
    arrivalAirport string,
    departureDate string,
    returnDate string,
    class string,
    airlineIDs []int,
    maxStops int,
    page int,
    limit int,
    sortBy string,
    sortOrder string,
    forUser bool,
) (*dto.RoundtripSearchResult, error) {
    log.Printf("Searching round-trip flights: from %s to %s, departure date: %s, return date: %s, class: %s",
        departureAirport, arrivalAirport, departureDate, returnDate, class)
    
    // Get outbound flights
    outboundFlights, err := r.SearchFlights(
        departureAirport,
        arrivalAirport,
        departureDate,
        class,
        airlineIDs,
        maxStops,
        page,
        limit,
        sortBy,
        sortOrder,
        forUser,
    )
    if err != nil {
        log.Printf("Error searching outbound flights: %v", err)
        return nil, fmt.Errorf("error searching outbound flights: %w", err)
    }

    // Get inbound flights
    inboundFlights, err := r.SearchFlights(
        arrivalAirport,
        departureAirport,
        returnDate,
        class,
        airlineIDs,
        maxStops,
        page,
        limit,
        sortBy,
        sortOrder,
        forUser,
    )
    if err != nil {
        log.Printf("Error searching inbound flights: %v", err)
        return nil, fmt.Errorf("error searching inbound flights: %w", err)
    }

    return &dto.RoundtripSearchResult{
        OutboundFlights: outboundFlights,
        InboundFlights:  inboundFlights,
    }, nil
}


func (r *flightRepository) UpdateFlightTime(flightID int, departureTime time.Time, arrivalTime time.Time) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    duration := arrivalTime.Sub(departureTime)
    durationMinutes := int(duration.Minutes())
    //check duration minutes
    if durationMinutes <= 0 {
    return fmt.Errorf("invalid flight duration: arrival time must be after departure time")
    }
    query := `
        UPDATE flights
        SET departure_time = $1, arrival_time = $2, duration_minutes = $3, updated_at = $4
        WHERE flight_id = $5
    `

    commandTag, err := r.db.Exec(ctx, query, departureTime, arrivalTime, durationMinutes, time.Now(), flightID)
    if err != nil {
        return fmt.Errorf("error updating flight times and duration: %w", err)
    }

    if commandTag.RowsAffected() == 0 {
        return fmt.Errorf("flight with ID %d not found", flightID)
    }

    log.Printf("Successfully updated flight %d: departure_time=%v, arrival_time=%v, duration=%d minutes", 
     flightID, departureTime, arrivalTime, durationMinutes)

    return nil
}