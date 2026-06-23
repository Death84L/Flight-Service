package interfaces

import (
	model "flight-service/internals/models"
)

type SeatRepository interface {
	GetSeatsByAirplaneID(airplaneID string) ([]model.Seat, error)
	GetSeatByAirplaneIDAndSeatNumber(airplaneID, seatNumber string) (*model.Seat, error)
	UpdateSeatStatus(seatID string, status string) error
}
