package interfaces

import (
	model "flight-service/internals/models"
)

type OrderRepository interface {
	CreateOrder(order model.Order) error
	GetOrderByFlightAndSeat(flightID, seatID string) (*model.Order, error)
	UpdateOrderStatus(orderID string, status string) error
}
