package interfaces

import (
	"context"

	model "flight-service/internals/models"
)

type FlightRepository interface {
	GetFlights(ctx context.Context, source, destination, date string, limit, offset int) ([]model.Flight, error)
	GetFlightByNumber(ctx context.Context, flightNumber string) (*model.Flight, error)
}
