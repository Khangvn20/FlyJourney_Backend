package service

import (
    "log"
    "time"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/dto"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/entity/error_code"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/model/request"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/model/response"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/port/repository"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/port/service"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/common/utils"
)
type flightService struct {
	flightRepo repository.FlightRepository
	
}
func NewFlightService(flightRepo repository.FlightRepository) service.FlightService {
	return &flightService{
		flightRepo: flightRepo,
	}
}
func (s *flightService) CreateFlightClasses(flightID int, req []request.FlightClassRequest) (*response.Response, error) {
    if len(req) == 0 {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "At least one flight class must be provided",
        }, nil
    }

    // Convert request to DTO
    flightClasses := make([]*dto.FlightClass, 0, len(req))
    for _, fcReq := range req {
        if fcReq.TotalSeats <= 0 {
            return &response.Response{
                Status:       false,
                ErrorCode:    error_code.InvalidRequest,
                ErrorMessage: "Total seats for each flight class must be greater than zero",
            }, nil
        }
        
        infantPrice := float64(0)
        if fcReq.BasePriceInfant != nil {
            infantPrice = *fcReq.BasePriceInfant
        }
        flightClasses = append(flightClasses, &dto.FlightClass{
            FlightID:         flightID,
            Class:            fcReq.Class,
            FareClassCode:    fcReq.FareClassCode,
            BasePrice:        fcReq.BasePrice,
            AvailableSeats:   fcReq.AvailableSeats,
            TotalSeats:       fcReq.TotalSeats,
            BasePriceChild:   fcReq.BasePriceChild,
            BasePriceInfant: infantPrice,
        })
    }

    createdClasses, err := s.flightRepo.CreateFlightClasses(flightID, flightClasses)
    if err != nil {
        log.Printf("Error creating flight classes: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Failed to create flight classes",
        }, err
    }

    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Flight classes created successfully",
        Data: map[string]interface{}{
            "flight_classes": createdClasses,
        },
    }, nil
}
func (s *flightService) CreateFlight(req *request.CreateFlightRequest) *response.Response {
		existingFlight, err := s.flightRepo.GetByFlightNumber(req.FlightNumber)
    if err != nil {
        log.Printf("Error checking flight number: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Failed to check flight number",
        }
    }
	    if existingFlight != nil {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Flight number already exists",
        }
    }
	  if len(req.FlightClasses) == 0 {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "At least one flight class must be provided",
        }
    }
    departureTime, err := utils.ParseTime(req.DepartureTime)
    if err != nil {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Invalid departure time format",
        }
        
    }
    arrivalTime, err := utils.ParseTime(req.ArrivalTime)
    if err != nil {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Invalid arrival time format",
        }
    }
    if arrivalTime.Before(departureTime) {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Arrival time must be after departure time",
        }
    }
    if departureTime.Before(time.Now()) {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Departure time must be in the future",
        }
    }
    totalSeats := 0
    for _, fc := range req.FlightClasses {
        if fc.TotalSeats <= 0 {
            return &response.Response{
                Status:       false,
                ErrorCode:    error_code.InvalidRequest,
                ErrorMessage: "Total seats for each flight class must be greater than zero",
            }
        }
        totalSeats += fc.TotalSeats
    }
	 flight := &dto.Flight{
        AirlineID:        req.AirlineID,
        FlightNumber:     req.FlightNumber,
        DepartureAirport: req.DepartureAirport,
        ArrivalAirport:   req.ArrivalAirport,
        DepartureTime:    departureTime,
        ArrivalTime:      arrivalTime,
        DurationMinutes:  req.DurationMinutes,
        StopsCount:       req.StopsCount,
        TaxAndFees:       req.TaxAndFees,
        ArrivalAiportCode:req.ArrivalAirportCode,
        DepartureAirportCode: req.DepartureCode,
        Currency:         req.Currency,
        TotalSeats:       totalSeats,
        Status:           req.Status,
        Distance:         req.Distance,
     
    }
	 createdFlight, err := s.flightRepo.CreateFlight(flight)
    if err != nil {
        log.Printf("Error creating flight: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: error_code.InternalErrMsg,
        }
    }
    flightClasses := make([]*dto.FlightClass, 0, len(req.FlightClasses))
    for _, fcReq := range req.FlightClasses {
         totalSeats += fcReq.TotalSeats
           infantPrice := float64(0)
        if fcReq.BasePriceInfant != nil {
            infantPrice = *fcReq.BasePriceInfant
        }
        
        flightClasses = append(flightClasses, &dto.FlightClass{
            FlightID:       createdFlight.FlightID, 
            Class:          fcReq.Class,
            BasePrice:      fcReq.BasePrice,
            AvailableSeats: fcReq.AvailableSeats,
            TotalSeats:     fcReq.TotalSeats,
            BasePriceChild: fcReq.BasePriceChild,
            BasePriceInfant: infantPrice,
            FareClassCode:  fcReq.FareClassCode,
        })
    }
    createdClasses, err := s.flightRepo.CreateFlightClasses(createdFlight.FlightID, flightClasses)
    if err != nil {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Error creating flight classes",
            Data:         nil,
        }
    }
    
	   return &response.Response{
        Status:       true,
        ErrorCode:    "",
        ErrorMessage: "",
        Data: map[string]interface{}{
            "flight":         createdFlight,
            "flight_classes": createdClasses, 
        },
    }      
    
}

func (s *flightService) UpdateFlight(flightID int, req *request.UpdateFlightRequest) *response.Response {
	    existingFlight, flightClasses, err := s.flightRepo.GetByID(flightID)

    if err != nil {
        log.Printf("Error getting flight: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Failed to get flight",
        }
    }
    if existingFlight == nil {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Not found flight with ID ",
        }
    }
	   if req.AirlineID != 0 {
        existingFlight.AirlineID = req.AirlineID
    }
    
    if req.FlightNumber != "" && req.FlightNumber != existingFlight.FlightNumber {
        flight, err := s.flightRepo.GetByFlightNumber(req.FlightNumber)
        if err != nil {
            log.Printf("Error checking flight number: %v", err)
            return &response.Response{
                Status:       false,
                ErrorCode:    error_code.InternalError,
                ErrorMessage: error_code.InternalErrMsg,
            }
        }
        if flight != nil {
            return &response.Response{
                Status:       false,
                ErrorCode:    error_code.InvalidRequest,
                ErrorMessage: "Flight number already exists",
            }
        }
        existingFlight.FlightNumber = req.FlightNumber
    }

    if req.DepartureAirport != "" {
        existingFlight.DepartureAirport = req.DepartureAirport
    }
    
    if req.ArrivalAirport != "" {
        existingFlight.ArrivalAirport = req.ArrivalAirport
    }

    if req.DepartureTime != "" {
        departureTime, err := utils.ParseTime(req.DepartureTime)
        if err != nil {
            return &response.Response{
                Status:       false,
                ErrorCode:    error_code.InvalidRequest,
                ErrorMessage: "Invalid departure time format",
            }
        }
        existingFlight.DepartureTime = departureTime
    }

    if req.ArrivalTime != "" {
        arrivalTime, err := utils.ParseTime(req.ArrivalTime)
        if err != nil {
            return &response.Response{
                Status:       false,
                ErrorCode:    error_code.InvalidRequest,
                ErrorMessage: "Invalid arrival time format",
            }
        }
        existingFlight.ArrivalTime = arrivalTime
    }

    if req.DurationMinutes != 0 {
        existingFlight.DurationMinutes = req.DurationMinutes
    }

    if req.StopsCount >= 0 {
        existingFlight.StopsCount = req.StopsCount
    }

    if req.TaxAndFees >= 0 {
        existingFlight.TaxAndFees = req.TaxAndFees
    }

    if req.Status != "" {
        existingFlight.Status = req.Status
    }

    if req.Distance != 0 {
        existingFlight.Distance = req.Distance
    }

    if req.Currency != "" {
        existingFlight.Currency = req.Currency
    }

    if req.DepartureAirportCode != "" {
        existingFlight.DepartureAirportCode = req.DepartureAirportCode
    }

    if req.ArrivalAirportCode != "" {
        existingFlight.ArrivalAiportCode = req.ArrivalAirportCode
    }

    if !existingFlight.ArrivalTime.IsZero() && !existingFlight.DepartureTime.IsZero() {
        if existingFlight.ArrivalTime.Before(existingFlight.DepartureTime) {
            return &response.Response{
                Status:       false,
                ErrorCode:    error_code.InvalidRequest,
                ErrorMessage: "Arrival time must be after departure time",
            }
        }
    }
   
    updatedFlight, err := s.flightRepo.Update(flightID, existingFlight)
    if err != nil {
        log.Printf("Error updating flight: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Failed to update flight",
        }
    }
    
    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Updated flight successfully",
        Data: map[string]interface{}{
            "flight":         updatedFlight,
            "flight_classes": flightClasses,
        },
    }
}
func (s *flightService) GetFlightByID(id int) *response.Response {
    flight, flightClasses, err := s.flightRepo.GetByID(id)
    if err != nil {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Failed to get flight",
            Data:         nil,
        }
    }

    if flight == nil {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Flight not found",
            Data:         nil,
        }
    }

    return &response.Response{
        Status:       true,
        ErrorCode:    "",
        ErrorMessage: "",
        Data: map[string]interface{}{
            "flight":         flight,
            "flight_classes": flightClasses,
        },
    }
}
func (s *flightService) GetAllFlights(page, limit int) *response.Response {
    flights, err := s.flightRepo.GetAll(page, limit)
    if err != nil {
        log.Printf("Error getting flights: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: error_code.InternalErrMsg,
        }
    }

    totalCount, err := s.flightRepo.Count()
    if err != nil {
        log.Printf("Error counting flights: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: error_code.InternalErrMsg,
        }
    }

    totalPages := (totalCount + limit - 1) / limit
    flightData := make([]map[string]interface{}, 0, len(flights))
    
    for _, flight := range flights {
        flightInfo := map[string]interface{}{
            "flight_id":               flight.FlightID,
            "airline_id":              flight.AirlineID,
            "flight_number":           flight.FlightNumber,
            "departure_airport":       flight.DepartureAirport,
            "arrival_airport":         flight.ArrivalAirport,
            "departure_time":          flight.DepartureTime,
            "arrival_time":            flight.ArrivalTime,
            "duration_minutes":        flight.DurationMinutes,
            "stops_count":             flight.StopsCount,
            "tax_and_fees":            flight.TaxAndFees,
            "total_seats":             flight.TotalSeats,
            "status":                  flight.Status,
            "distance":                flight.Distance,
            "departure_airport_code":  flight.DepartureAirportCode,
            "arrival_airport_code":    flight.ArrivalAiportCode,
            "currency":                flight.Currency,
            "created_at":              flight.CreatedAt,
            "updated_at":              flight.UpdatedAt,
        }
      

        flightData = append(flightData, flightInfo)
    }

    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Successfully retrieved flights",
        Data: map[string]interface{}{
            "flights":     flightData,
            "total_count": totalCount,
            "page":        page,
            "limit":       limit,
            "total_pages": totalPages,
        },
    }
}
func (s *flightService) SearchFlights(req *request.FlightSearchRequest) *response.Response {
    // Validate and parse departure date
    _, err := utils.ParseTime(req.DepartureDate)
    if err != nil {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Invalid departure date format. Supported formats: dd/mm/yyyy, yyyy-mm-dd, dd/mm/yyyy HH:mm",
        }
    }

    // Validate and parse arrival date if provided
    if req.ArrivalDate != "" {
        _, err := utils.ParseTime(req.ArrivalDate)
        if err != nil {
            return &response.Response{
                Status:       false,
                ErrorCode:    error_code.InvalidRequest,
                ErrorMessage: "Invalid arrival date format. Supported formats: dd/mm/yyyy, yyyy-mm-dd, dd/mm/yyyy HH:mm",
            }
        }
    }

    var airlineIDs []int
    if req.AirlineIDs != nil {
        airlineIDs = req.AirlineIDs
    }
    maxStops := -1
    if req.MaxStops >= 0 {
        maxStops = req.MaxStops
    }
    
    // Set default pagination values
    page := 1
    if req.Page > 0 {
        page = req.Page
    }
    
    limit := 10
    if req.Limit > 0 {
        limit = req.Limit
    }
    
    // Define sort parameters with defaults
    sortBy := "departure_time"
    if req.SortBy != "" {
        sortBy = req.SortBy
    }
    
    sortOrder := "ASC"
    if req.SortOrder != "" {
        sortOrder = req.SortOrder
    }
    
    flights, err := s.flightRepo.SearchFlights(
        req.DepartureAirportCode,
        req.ArrivalAirportCode,
        req.DepartureDate,
        req.FlightClass,
        airlineIDs,
        maxStops,
        page,
        limit,
        sortBy,
        sortOrder,
        false,
    )
    
    if err != nil {
        log.Printf("Error searching flights: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: error_code.InternalErrMsg,
        }
    }
    
    totalCount, err := s.flightRepo.CountBySearch(
        req.DepartureAirportCode,
        req.ArrivalAirportCode,
        req.DepartureDate,
        false, // For admin search, we don't filter by status
    )
    
    if err != nil {
        log.Printf("Error counting flights: %v", err)
        totalCount = len(flights)
    }
    
    totalPages := (totalCount + limit - 1) / limit
    
    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Successfully searched flights",
        Data: map[string]interface{}{
            "search_results":    flights, // Changed from "flights" to "search_results"
            "total_count":       totalCount,
            "total_pages":       totalPages,
            "page":              page,
            "limit":             limit,
            "departure_airport": req.DepartureAirportCode,
            "arrival_airport":   req.ArrivalAirportCode,
            "departure_date":    req.DepartureDate,
            "arrival_date":      req.ArrivalDate,
            "flight_class":      req.FlightClass,
            "sort_by":           sortBy,
            "sort_order":        sortOrder,
        },
    }
}
func (s *flightService) SearchFlightsForUser(req *request.FlightSearchRequest) *response.Response {
    // Validate and parse departure date
    _, err := utils.ParseTime(req.DepartureDate)
    if err != nil {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Invalid departure date format. Supported formats: dd/mm/yyyy, yyyy-mm-dd, dd/mm/yyyy HH:mm",
        }
    }


    // Validate and parse arrival date if provided
    if req.ArrivalDate != "" {
        _, err := utils.ParseTime(req.ArrivalDate)
        if err != nil {
            return &response.Response{
                Status:       false,
                ErrorCode:    error_code.InvalidRequest,
                ErrorMessage: "Invalid arrival date format. Supported formats: dd/mm/yyyy, yyyy-mm-dd, dd/mm/yyyy HH:mm",
            }
        }
    }

    var airlineIDs []int
    if req.AirlineIDs != nil {
        airlineIDs = req.AirlineIDs
    }

    maxStops := -1
    if req.MaxStops >= 0 {
        maxStops = req.MaxStops
    }

    page := 1
    if req.Page > 0 {
        page = req.Page
    }

    limit := 10
    if req.Limit > 0 {
        limit = req.Limit
    }

    sortBy := "departure_time"
    if req.SortBy != "" {
        sortBy = req.SortBy
    }

    sortOrder := "ASC"
    if req.SortOrder != "" {
        sortOrder = req.SortOrder
    }

    flights, err := s.flightRepo.SearchFlights(
        req.DepartureAirportCode,
        req.ArrivalAirportCode,
        req.DepartureDate,
        req.FlightClass,
        airlineIDs,
        maxStops,
        page,
        limit,
        sortBy,
        sortOrder,
        true, // forUser = true
    )

    if err != nil {
        log.Printf("Error searching flights for user: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Failed to search flights",
        }
    }

    for _, flight := range flights {
        passengers := req.Passengers
        
        // Always calculate adult cost (minimum 1 adult required)
        totalAdultCost := flight.Pricing.TotalPrices.Adult * float64(passengers.Adults)
        
        var totalChildCost, totalInfantCost float64
        
        if passengers.Children > 0 {
            totalChildCost = flight.Pricing.TotalPrices.Child * float64(passengers.Children)
        } else {
            // Clear child data if not needed
            flight.Pricing.BasePrices.Child = 0
            flight.Pricing.TotalPrices.Child = 0
            flight.Pricing.Taxes.Child = 0
        }
    
        if passengers.Infants > 0 {
            totalInfantCost = flight.Pricing.TotalPrices.Infant * float64(passengers.Infants)
        } else {

            flight.Pricing.BasePrices.Infant = 0
            flight.Pricing.TotalPrices.Infant = 0
            flight.Pricing.Taxes.Infant = 0
        }
        
        flight.Pricing.GrandTotal = totalAdultCost + totalChildCost + totalInfantCost
    }

    totalCount, err := s.flightRepo.CountBySearch(
        req.DepartureAirportCode,
        req.ArrivalAirportCode,
        req.DepartureDate,
        true, // forUser = true
    )
    
    if err != nil {
        log.Printf("Error counting flights for user: %v", err)
        totalCount = len(flights)
    }

    totalPages := (totalCount + limit - 1) / limit

    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Successfully searched flights for user",
        Data: map[string]interface{}{
            "search_results":    flights, // Using FlightSearchResult struct
            "total_count":       totalCount,
            "total_pages":       totalPages,
            "page":              page,
            "limit":             limit,
            "departure_airport": req.DepartureAirportCode,
            "arrival_airport":   req.ArrivalAirportCode,
            "departure_date":    req.DepartureDate,
            "arrival_date":      req.ArrivalDate,
            "flight_class":      req.FlightClass,
             "passengers":        req.Passengers,
            "sort_by":           sortBy,
            "sort_order":        sortOrder,
        },
    }
}
func (s *flightService) UpdateFlightStatus(flightID int, req *request.UpdateFlightStatusRequest) *response.Response {
	err := s.flightRepo.UpdateStatus(flightID, req.Status)
    if err != nil {
		if err.Error() == "flight not found" {
			return &response.Response{
				Status:       false,
				ErrorCode:    error_code.InvalidRequest,
				ErrorMessage: "Flight not found with ID " ,
			}
		}
		log.Printf("Error updating flight status: %v", err)
		return &response.Response{
			Status:       false,
			ErrorCode:    error_code.InternalError,
			ErrorMessage: "Failed to update flight status",
		}
	}
	return &response.Response{
		Status:       true,
		ErrorCode:    error_code.Success,
		ErrorMessage: "Flight status updated successfully",
	}
}
func (s *flightService) GetFlightByAirline(airlineID int, page, limit int) *response.Response {
    flights, err := s.flightRepo.GetByAirline(airlineID, page, limit)
    if err != nil {
        log.Printf("Error getting flights by airline: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:   error_code.InternalError,
            ErrorMessage: "Failed to get flights by airline",
        }
    }
    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Successfully retrieved flights by airline",
        Data: flights,
    }
}
func (s *flightService) GetFareCLassCode(flightID int) *response.Response {
    fareClasses, err := s.flightRepo.GetFareClassesByFlightID(flightID)
    if err != nil {
        log.Printf("Error getting fare class code: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Failed to get fare class code",
        }
    }
    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Successfully retrieved fare class code",
        Data:         fareClasses,
    }
}
func (s *flightService) GetFlightsByDate(req *request.GetFlightsByDateRequest) *response.Response {
      parsedDate, err := utils.ParseTime(req.Date)
       if err != nil {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Invalid date format. Supported formats: dd/mm/yyyy, yyyy-mm-dd, dd/mm/yyyy HH:mm",
        }
    }
    flights, err := s.flightRepo.GetFlightsByDate(parsedDate, req.Page, req.Limit)
    if err != nil {
        log.Printf("Error getting flights by date: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Failed to get flights by date",
        }
    }
    if req.Status != "" {
        filteredFlights := make([]*dto.Flight, 0)
        for _, flight := range flights {
            if flight.Status == req.Status {
                filteredFlights = append(filteredFlights, flight)
            }
        }
        flights = filteredFlights
    }
    totalCount, err := s.flightRepo.CountByDate(parsedDate)
    if err != nil {
        log.Printf("Error counting flights by date: %v", err)
        totalCount = len(flights)
    }
    if req.Limit <= 0 {
        req.Limit = 10
    }
    if req.Page < 1 {
        req.Page = 1
    }
    totalPages := (totalCount + req.Limit - 1) / req.Limit

    // Transform data cho response
    flightData := make([]map[string]interface{}, 0, len(flights))
    for _, flight := range flights {
        flightInfo := map[string]interface{}{
            "flight_id":               flight.FlightID,
            "airline_id":              flight.AirlineID,
            "flight_number":           flight.FlightNumber,
            "departure_airport":       flight.DepartureAirport,
            "arrival_airport":         flight.ArrivalAirport,
            "departure_time":          flight.DepartureTime,
            "arrival_time":            flight.ArrivalTime,
            "duration_minutes":        flight.DurationMinutes,
            "stops_count":             flight.StopsCount,
            "tax_and_fees":            flight.TaxAndFees,
            "total_seats":             flight.TotalSeats,
            "status":                  flight.Status,
            "distance":                flight.Distance,
            "departure_airport_code":  flight.DepartureAirportCode,
            "arrival_airport_code":    flight.ArrivalAiportCode,
            "currency":                flight.Currency,
            "created_at":              flight.CreatedAt,
            "updated_at":              flight.UpdatedAt,
            "flight_classes":          flight.FlightClasses,
        }
        flightData = append(flightData, flightInfo)
    }

    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Successfully retrieved flights by date",
        Data: map[string]interface{}{
            "flights":     flightData,
            "date":        req.Date,
            "total_count": totalCount,
            "page":        req.Page,
            "limit":       req.Limit,
            "total_pages": totalPages,
            "status_filter": req.Status,
        },
    }
}
func (s *flightService) GetFlightsByStatus(status string, page, limit int) *response.Response {
    flights, err := s.flightRepo.GetByStatus(status, page, limit)
    if err != nil {
        log.Printf("Error getting flights by status: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Failed to get flights by status",
        }
    }
    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Successfully retrieved flights by status",
        Data:         flights,
    }
}
func (s *flightService) SearchRoundtripFlights(req *request.RoundtripFlightSearchRequest) *response.Response {
    var airlineIDs []int
    if req.AirlineIDs != nil {
        airlineIDs = req.AirlineIDs
    }
    maxStops := -1
    if req.MaxStops >= 0 {
        maxStops = req.MaxStops
    }
    page := 1
    if req.Page > 0 {
        page = req.Page
    }
    limit := 10
    if req.Limit > 0 {
        limit = req.Limit
    }
    sortBy := "departure_time"
    if req.SortBy != "" {
        sortBy = req.SortBy
    }
    sortOrder := "ASC"
    if req.SortOrder != "" {
        sortOrder = req.SortOrder
    }

    // Tìm chuyến bay chiều đi
    outboundFlights, err := s.flightRepo.SearchFlights(
        req.DepartureAirportCode,
        req.DepartureAirportCode,
        req.DepartureDate,
        req.FlightClass,
        airlineIDs,
        maxStops,
        page,
        limit,
        sortBy,
        sortOrder,
        false, // Admin: lấy tất cả status
    )
    if err != nil {
        log.Printf("Error searching outbound flights: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: error_code.InternalErrMsg,
        }
    }

    // Tìm chuyến bay chiều về
    inboundFlights, err := s.flightRepo.SearchFlights(
        req.ArrivalAirportCode,
        req.DepartureAirportCode,
       req.ReturnDate,
        req.FlightClass,
        airlineIDs,
        maxStops,
        page,
        limit,
        sortBy,
        sortOrder,
        false, // Admin: lấy tất cả status
    )
    //Caculate total count for outbound and inbound flights
               for _, flight := range outboundFlights {
    passengers := req.Passengers
    totalAdultCost := flight.Pricing.TotalPrices.Adult * float64(passengers.Adults)
    var totalChildCost, totalInfantCost float64

    if passengers.Children > 0 {
        totalChildCost = flight.Pricing.TotalPrices.Child * float64(passengers.Children)
    } else {
        flight.Pricing.BasePrices.Child = 0
        flight.Pricing.TotalPrices.Child = 0
        flight.Pricing.Taxes.Child = 0
    }

    if passengers.Infants > 0 {
        totalInfantCost = flight.Pricing.TotalPrices.Infant * float64(passengers.Infants)
    } else {
        flight.Pricing.BasePrices.Infant = 0
        flight.Pricing.TotalPrices.Infant = 0
        flight.Pricing.Taxes.Infant = 0
    }

    flight.Pricing.GrandTotal = totalAdultCost + totalChildCost + totalInfantCost
}

    for _, flight := range inboundFlights {
    passengers := req.Passengers
    totalAdultCost := flight.Pricing.TotalPrices.Adult * float64(passengers.Adults)
    var totalChildCost, totalInfantCost float64

    if passengers.Children > 0 {
        totalChildCost = flight.Pricing.TotalPrices.Child * float64(passengers.Children)
    } else {
        flight.Pricing.BasePrices.Child = 0
        flight.Pricing.TotalPrices.Child = 0
        flight.Pricing.Taxes.Child = 0
    }

    if passengers.Infants > 0 {
        totalInfantCost = flight.Pricing.TotalPrices.Infant * float64(passengers.Infants)
    } else {
        flight.Pricing.BasePrices.Infant = 0
        flight.Pricing.TotalPrices.Infant = 0
        flight.Pricing.Taxes.Infant = 0
    }

    flight.Pricing.GrandTotal = totalAdultCost + totalChildCost + totalInfantCost
}
    if err != nil {
        log.Printf("Error searching inbound flights: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: error_code.InternalErrMsg,
        }
    }

    // Đếm tổng số chuyến bay chiều đi
    outboundCount, err := s.flightRepo.CountBySearch(
        req.DepartureAirportCode,
        req.ArrivalAirportCode,
        req.DepartureDate,
        false, // Admin: lấy tất cả status
    )
    if err != nil {
        outboundCount = len(outboundFlights)
    }

    inboundCount, err := s.flightRepo.CountBySearch(
        req.ArrivalAirportCode,
        req.DepartureAirportCode,
        req.ReturnDate,
        false, // Admin: lấy tất cả status
    )
    if err != nil {
        inboundCount = len(inboundFlights)
    }

    outboundTotalPages := (outboundCount + limit - 1) / limit
    inboundTotalPages := (inboundCount + limit - 1) / limit

    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Successfully searched roundtrip flights",
        Data: map[string]interface{}{
            "outbound_flights":      outboundFlights,
            "inbound_flights":       inboundFlights,
            "outbound_total_count":  outboundCount,
            "inbound_total_count":   inboundCount,
            "outbound_total_pages":  outboundTotalPages,
            "inbound_total_pages":   inboundTotalPages,
            "page":                  page,
            "limit":                 limit,
            "departure_airport":     req.DepartureAirportCode,
            "arrival_airport":       req.ArrivalAirportCode,
            "departure_date":        req.DepartureDate,
            "return_date":           req.ReturnDate,
            "flight_class":          req.FlightClass,
            "passenger_count":       req.Passengers,
            "sort_by":               sortBy,
            "sort_order":            sortOrder,
        },
    }
}
func (s *flightService) GetFlightByIDForUser(flightID int) *response.Response {
    flight, flightClasses, err := s.flightRepo.GetByID(flightID)                            // ✅ Debug
    
    if err != nil {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError, 
            ErrorMessage: "Fail to get info for user",
            Data:         nil,
        }
    }
    if flight == nil {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Flight not found",
            Data:         nil,
        }
    }
    if flight.Status != "scheduled" && flight.Status != "boarding" {
        return &response.Response{  
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Flight is not available for booking",
        }
    }
    userFlight := &dto.UserFlightDetail{
        FlightID:         flight.FlightID,
        AirlineID:        flight.AirlineID,
        FlightNumber:     flight.FlightNumber,
        DepartureAirport: flight.DepartureAirport,
        ArrivalAirport:   flight.ArrivalAirport,
        DepartureTime:    flight.DepartureTime,
        ArrivalTime:      flight.ArrivalTime,
        DurationMinutes:  flight.DurationMinutes,
        StopsCount:       flight.StopsCount,
        TaxAndFees:       flight.TaxAndFees,
        Distance:         flight.Distance,
        FlightClasses:    make([]*dto.UserFlightClass, 0, len(flightClasses)),
    }
        for _, fc := range flightClasses {
        userFlightClass := &dto.UserFlightClass{
            FlightClassID:    fc.FlightClassID,
            Class:            fc.Class,
            BasePrice:        fc.BasePrice,
            AvailableSeats:   fc.AvailableSeats,
            BasePriceChild:   fc.BasePriceChild,
            BasePriceInfant:  fc.BasePriceInfant,
            FareClassCode:    fc.FareClassCode,
            FareClassDetails: fc.FareClassDetails, 
        }
        userFlight.FlightClasses = append(userFlight.FlightClasses, userFlightClass)

    }
    
    return &response.Response{
        Status:      true,
        ErrorCode:  error_code.Success,
        ErrorMessage: "Successfully retrieved flight details for user",
         Data:         userFlight,
    }
}
func (s *flightService) GetFlightByIDForAdmin(flightID int) *response.Response {
    flight, flightClasses, err := s.flightRepo.GetByID(flightID)
    if err != nil {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Fail to get info for admin",
            Data:         nil,
        }
    }
    if flight == nil {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Flight not found",
            Data:         nil,
        }
    }
    adminFlight := &dto.AdminFlightDetail{
        FlightID:         flight.FlightID,
        AirlineID:        flight.AirlineID,
        FlightNumber:     flight.FlightNumber,
        DepartureAirport: flight.DepartureAirport,
        ArrivalAirport:   flight.ArrivalAirport,
        DepartureTime:    flight.DepartureTime,
        ArrivalTime:      flight.ArrivalTime,
        DurationMinutes:  flight.DurationMinutes,
        StopsCount:       flight.StopsCount,
        TaxAndFees:       flight.TaxAndFees,
        TotalSeats:       flight.TotalSeats,
        Status:           flight.Status,
        Distance:         flight.Distance,
        CreatedAt:        flight.CreatedAt,
        UpdatedAt:        flight.UpdatedAt,
        FlightClasses:    make([]*dto.AdminFlightClass, 0, len(flightClasses)),
    }
    for _, fc := range flightClasses {
        adminFlightClass := &dto.AdminFlightClass{
            FlightClassID:    fc.FlightClassID,
            FlightID:         fc.FlightID,
            Class:            fc.Class,
            BasePrice:        fc.BasePrice,
            AvailableSeats:   fc.AvailableSeats,
            BasePriceChild:   fc.BasePriceChild,
            BasePriceInfant:  fc.BasePriceInfant,
            FareClassCode:    fc.FareClassCode,
            TotalSeats:       fc.TotalSeats,
            FareClassDetails: fc.FareClassDetails,
            CreatedAt:        fc.CreatedAt,
            UpdatedAt:        fc.UpdatedAt,
        }
        adminFlight.FlightClasses = append(adminFlight.FlightClasses, adminFlightClass)
    }
    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Successfully retrieved flight details",
        Data:         adminFlight,
    }
}

func (s *flightService) SearchRoundtripFlightsForUser(req *request.RoundtripFlightSearchRequest) *response.Response {
    var airlineIDs []int
    if req.AirlineIDs != nil {
        airlineIDs = req.AirlineIDs
    }
    maxStops := -1
    if req.MaxStops >= 0 {
        maxStops = req.MaxStops
    }
    page :=1
    if req.Page > 0 {
        page = req.Page
    }
    limit :=10
    if req.Limit > 0 {
        limit = req.Limit
    }
    sortBy := "departure_time"
    if req.SortBy != "" {
        sortBy = req.SortBy
    }
    sortOrder := "ASC"
    if req.SortOrder != "" {
        sortOrder = req.SortOrder
    }
    
    outboundFlights, err := s.flightRepo.SearchFlights(
        req.DepartureAirportCode,
        req.ArrivalAirportCode,
        req.DepartureDate,
        req.FlightClass,
        airlineIDs,
        maxStops,
        page,
        limit,
        sortBy,
        sortOrder,
        true,
    )
    if err !=nil {
        log.Printf("Error searching outbound flights: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Error searching outbound flights",
        }
    }
    inboundFlights, err := s.flightRepo.SearchFlights(
        req.ArrivalAirportCode,
        req.DepartureAirportCode,
        req.ReturnDate,
        req.FlightClass,
        airlineIDs,
        maxStops,
        page,
        limit,
        sortBy,
        sortOrder,
        true,
    )

    //Cacutalate total count for outbound and inbound flights
                  for _, flight := range outboundFlights {
    passengers := req.Passengers
    totalAdultCost := flight.Pricing.TotalPrices.Adult * float64(passengers.Adults)
    var totalChildCost, totalInfantCost float64

    if passengers.Children > 0 {
        totalChildCost = flight.Pricing.TotalPrices.Child * float64(passengers.Children)
    } else {
        flight.Pricing.BasePrices.Child = 0
        flight.Pricing.TotalPrices.Child = 0
        flight.Pricing.Taxes.Child = 0
    }

    if passengers.Infants > 0 {
        totalInfantCost = flight.Pricing.TotalPrices.Infant * float64(passengers.Infants)
    } else {
        flight.Pricing.BasePrices.Infant = 0
        flight.Pricing.TotalPrices.Infant = 0
        flight.Pricing.Taxes.Infant = 0
    }

    flight.Pricing.GrandTotal = totalAdultCost + totalChildCost + totalInfantCost
}

    for _, flight := range inboundFlights {
    passengers := req.Passengers
    totalAdultCost := flight.Pricing.TotalPrices.Adult * float64(passengers.Adults)
    var totalChildCost, totalInfantCost float64

    if passengers.Children > 0 {
        totalChildCost = flight.Pricing.TotalPrices.Child * float64(passengers.Children)
    } else {
        flight.Pricing.BasePrices.Child = 0
        flight.Pricing.TotalPrices.Child = 0
        flight.Pricing.Taxes.Child = 0
    }

    if passengers.Infants > 0 {
        totalInfantCost = flight.Pricing.TotalPrices.Infant * float64(passengers.Infants)
    } else {
        flight.Pricing.BasePrices.Infant = 0
        flight.Pricing.TotalPrices.Infant = 0
        flight.Pricing.Taxes.Infant = 0
    }

    flight.Pricing.GrandTotal = totalAdultCost + totalChildCost + totalInfantCost
}
    if err != nil {
        log.Printf("Error searching inbound flights: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Error searching inbound flights",
        }
    }
     outboundCount, err := s.flightRepo.CountBySearch(
        req.DepartureAirportCode,
        req.ArrivalAirportCode,
        req.DepartureDate,
        false,
    )
    if err != nil {
        outboundCount = len(outboundFlights)
    }

    inboundCount, err := s.flightRepo.CountBySearch(
        req.ArrivalAirportCode,
        req.DepartureAirportCode,
        req.ReturnDate,
        false,
    )
    if err != nil {
        inboundCount = len(inboundFlights)
    }

    outboundTotalPages := (outboundCount + limit - 1) / limit
    inboundTotalPages := (inboundCount + limit - 1) / limit

    roundtripResult := &dto.RoundtripSearchResult{
        OutboundFlights: outboundFlights,
        InboundFlights:  inboundFlights,
    }
    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Successfully searched roundtrip flights",
        Data: map[string]interface{}{
            "search_results":        roundtripResult,
            "outbound_total_count":    outboundCount,
            "inbound_total_count":     inboundCount,
            "outbound_total_pages":    outboundTotalPages,
            "inbound_total_pages":     inboundTotalPages,
            "page":                    page,
            "limit":                   limit,
            "departure_airport":       req.DepartureAirportCode,
            "arrival_airport":         req.ArrivalAirportCode,
            "departure_date":          req.DepartureDate,
            "return_date":             req.ReturnDate,
            "flight_class":            req.FlightClass,
            "passenger_count":         req.Passengers,
            "sort_by":                 sortBy,
            "sort_order":              sortOrder,
        },
    }
}