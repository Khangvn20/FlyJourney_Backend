package service

import (
    "log"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/dto"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/entity/error_code"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/model/request"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/model/response"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/port/repository"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/port/service"
)
type flightService struct {
	flightRepo repository.FlightRepository
	
}
func NewFlightService(flightRepo repository.FlightRepository) service.FlightService {
	return &flightService{
		flightRepo: flightRepo,
	}
}
func (s *flightService) CreateFlight(req *request.CreateFlightRequest) *response.Response {
		existingFlight, err := s.flightRepo.GetByFlightNumber(req.FlightNumber)
    if err != nil {
        log.Printf("Error checking flight number: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Flight number already exists",
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
        AircraftID:       req.AircraftID,
        FlightNumber:     req.FlightNumber,
        DepartureAirport: req.DepartureAirport,
        ArrivalAirport:   req.ArrivalAirport,
        DepartureTime:    req.DepartureTime,
        ArrivalTime:      req.ArrivalTime,
        DurationMinutes:  req.DurationMinutes,
        StopsCount:       req.StopsCount,
        TaxAndFees:       req.TaxAndFees,
        TotalSeats:       totalSeats,
        Status:           req.Status,
        Gate:             req.Gate,
        Terminal:         req.Terminal,
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
        flightClasses = append(flightClasses, &dto.FlightClass{
            FlightID:       createdFlight.FlightID, 
            Class:          fcReq.Class,
            BasePrice:      fcReq.BasePrice,
            AvailableSeats: fcReq.AvailableSeats,
            TotalSeats:     fcReq.TotalSeats,
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
            ErrorMessage: error_code.InternalErrMsg,
        }
    }
    if existingFlight == nil {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Not found flight with ID ",
        }
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
                ErrorMessage: "Flight number đã tồn tại",
            }
        }
    }
	 updatedFlight, err := s.flightRepo.Update(flightID, existingFlight)
    if err != nil {
        log.Printf("Error updating flight: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: error_code.InternalErrMsg,
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
            ErrorMessage: err.Error(),
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
func (s *flightService) GetAllFlights(page , limit int) *response.Response {
	flights ,err := s.flightRepo.GetAll(page, limit)
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
	return &response.Response{
		Status:       true,
		ErrorCode:    error_code.Success,
		ErrorMessage: "Successfully retrieved flights",
		  Data: map[string]interface{}{
            "flights":     flights,
            "total_count": totalCount,
            "page":        page,
            "limit":       limit,
            "total_pages": totalPages,
        },
	}
}
func (s *flightService) SearchFlights(req *request.FlightSearchRequest) *response.Response {
    // Initialize default values for optional parameters
    var airlineIDs []int
    if req.AirlineIDs != nil {
        airlineIDs = req.AirlineIDs
    }
    
    // Default stops value if not specified
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
    
    // Call repository with all required parameters
    flights, err := s.flightRepo.SearchFlights(
        req.DepartureAirport,
        req.ArrivalAirport,
        req.DepartureDate,
        req.FlightClass,
        airlineIDs,
        maxStops,
        page,
        limit,
        sortBy,
        sortOrder,
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
        req.DepartureAirport,
        req.ArrivalAirport,
        req.DepartureDate,
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
            "flights":           flights,
            "total_count":       totalCount,
            "total_pages":       totalPages,
            "page":              page,
            "limit":             limit,
            "departure_airport": req.DepartureAirport,
            "arrival_airport":   req.ArrivalAirport,
            "departure_date":    req.DepartureDate,
            "flight_class":      req.FlightClass,
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
			ErrorMessage: error_code.InternalErrMsg,
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
            ErrorMessage: error_code.InternalErrMsg,
        }
    }
    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Successfully retrieved flights by airline",
        Data: flights,
    }
}

func (s *flightService) GetFlightsByStatus(status string, page, limit int) *response.Response {
    flights, err := s.flightRepo.GetByStatus(status, page, limit)
    if err != nil {
        log.Printf("Error getting flights by status: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: error_code.InternalErrMsg,
        }
    }
    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Successfully retrieved flights by status",
        Data:         flights,
    }
}
