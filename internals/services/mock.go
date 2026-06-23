package services

import (
	"context"

	model "flight-service/internals/models"
)

/*
========================
FLIGHT MOCK
========================
*/

type mockFlightRepo struct {
	flight  *model.Flight
	flights []model.Flight
	err     error
}

func (m *mockFlightRepo) GetFlights(
	ctx context.Context,
	source, destination, date string,
	limit, offset int,
) ([]model.Flight, error) {
	return m.flights, m.err
}

func (m *mockFlightRepo) GetFlightByNumber(
	ctx context.Context,
	flightNumber string,
) (*model.Flight, error) {
	return m.flight, m.err
}

/*
========================
SEAT MOCK
========================
*/

type mockSeatRepo struct {
	seat  *model.Seat
	seats []model.Seat
	err   error

	updateStatusErr error // control UpdateSeatStatus independently
}

func (m *mockSeatRepo) GetSeatsByAirplaneID(airplaneID string) ([]model.Seat, error) {
	return m.seats, m.err
}

func (m *mockSeatRepo) GetSeatByAirplaneIDAndSeatNumber(
	airplaneID, seatNumber string,
) (*model.Seat, error) {
	return m.seat, m.err
}

func (m *mockSeatRepo) UpdateSeatStatus(seatID string, status string) error {
	return m.updateStatusErr
}

/*
========================
ORDER MOCK
========================
*/

type mockOrderRepo struct {
	existingOrder *model.Order
	createErr     error
	updateErr     error
}

func (m *mockOrderRepo) CreateOrder(order model.Order) error {
	return m.createErr
}

func (m *mockOrderRepo) GetOrderByFlightAndSeat(flightID, seatID string) (*model.Order, error) {
	return m.existingOrder, nil
}

func (m *mockOrderRepo) UpdateOrderStatus(orderID string, status string) error {
	return m.updateErr
}
