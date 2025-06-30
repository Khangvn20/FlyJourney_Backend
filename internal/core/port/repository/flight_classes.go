package repository
import (
	"context"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/dto"
)
type FlightClassesRepository interface {
	CreateFlightClass(ctx context.Context, flightClass *dto.FlightClass) (*dto.FlightClass, error)
	GetByID (id int) (*dto.FlightClass, error)
    GetByFlightID (flightID int) ([]*dto.FlightClass, error)
	Update (id int, flightClass *dto.FlightClass) (*dto.FlightClass, error)
	GetByFlightAndClass(flightID int, class string) (*dto.FlightClass, error)
	UpdateAvailableSeats(id int , Seats int) error
    CreateBulk(flightClasses []*dto.FlightClass) ([]*dto.FlightClass, error)
}