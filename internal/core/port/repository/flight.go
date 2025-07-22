package repository
import (
	"time"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/dto"
)	
type FlightRepository interface {
	 CreateFlight(flight *dto.Flight) (*dto.Flight, error)
	  GetByID(id int) (*dto.Flight, []*dto.FlightClass, error)
	 GetAll(page, limit int) ([]*dto.Flight, error)
	 Update(id int, flight *dto.Flight) (*dto.Flight, error)
	 CreateFlightClasses(flightID int, classes []*dto.FlightClass) ([]*dto.FlightClass, error)
	 //Specialize query methods
	 GetByFlightNumber(flightNumber string) (*dto.Flight, error)
     GetByRoute(departureAirport, arrivalAirport string, date time.Time) ([]*dto.Flight, error)
     GetByAirline(airlineID int, page, limit int) ([]*dto.Flight, error)
     GetByStatus(status string, page, limit int) ([]*dto.Flight, error)
	//Search methods
	 SearchFlights(
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
    ) ([]*dto.FlightSearchResult, error)
     SearchRoundtripFlights(
        departureAirport string,
        arrivalAirport string,
        departureDate string,
        returnDate   string,
        flightClass string,
        airlineIDs []int,
        maxStops int,
        page int,
        limit int,
        sortBy string,
        sortOrder string,
    ) (*dto.RoundtripSearchResult, error)
		//Metadata methods
		 Count() (int, error)
    CountBySearch(departureAirport, arrivalAirport string, departureDate string,forUser bool, ) (int, error)
    // Status updates
    UpdateStatus(id int, status string) error

}