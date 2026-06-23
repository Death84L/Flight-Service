package model

import "time"

type SeatStatus string

const (
	SeatAvailable SeatStatus = "AVAILABLE"
	SeatReserved  SeatStatus = "RESERVED"
	SeatBooked    SeatStatus = "BOOKED"
)

type Seat struct {
	ID string `json:"id"`

	AirplaneID string `json:"airplane_id"`

	SeatNumber string `json:"seat_number"`

	Price float64 `json:"price"`

	Status SeatStatus `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
