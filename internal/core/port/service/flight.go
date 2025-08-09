package service

import (
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/model/request"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/model/response"
)
type FlightService interface {
	CreateFlightClasses(flightID int, classes []request.FlightClassRequest) (*response.Response, error)
	CreateFlight(request *request.CreateFlightRequest) *response.Response
	GetFlightByID(flightID int) *response.Response
	GetAllFlights(page, limit int) *response.Response
	UpdateFlight(flightID int, request *request.UpdateFlightRequest) *response.Response
	UpdateFlightStatus(flightID int ,req *request.UpdateFlightStatusRequest) *response.Response
	SearchFlights(req *request.FlightSearchRequest) *response.Response
	GetFlightsByDate(req *request.GetFlightsByDateRequest) *response.Response
	SearchFlightsForUser(req *request.FlightSearchRequest) *response.Response
	GetFareCLassCode(flightID int) *response.Response
	GetFlightByAirline(airlineID int, page, limit int) *response.Response
	GetFlightsByStatus(status string, page, limit int) *response.Response
	SearchRoundtripFlights(req *request.RoundtripFlightSearchRequest) *response.Response
	GetFlightByIDForUser(flightID int) *response.Response 
    GetFlightByIDForAdmin(flightID int) *response.Response
	SearchRoundtripFlightsForUser(req *request.RoundtripFlightSearchRequest) *response.Response
}