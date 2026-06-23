package model

import "time"

type OrderStatus string

const (
	OrderPending   OrderStatus = "PENDING"
	OrderConfirmed OrderStatus = "CONFIRMED"
	OrderFailed    OrderStatus = "FAILED"
	OrderCancelled OrderStatus = "CANCELLED"
)

type Order struct {
	ID string `json:"id"`

	FlightID     string `json:"flight_id"`
	FlightNumber string `json:"flight_number"`

	SeatID     string `json:"seat_id"`
	SeatNumber string `json:"seat_number"`

	Price float64 `json:"price"`

	Status OrderStatus `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
