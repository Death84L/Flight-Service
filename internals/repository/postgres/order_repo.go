package postgres

import (
	"database/sql"
	model "flight-service/internals/models"
)

type OrderRepo struct {
	DB *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{DB: db}
}

func (r *OrderRepo) CreateOrder(order model.Order) error {

	query := `
	INSERT INTO orders (id, flight_id, seat_id, price, status)
	VALUES ($1,$2,$3,$4,$5)
	`

	_, err := r.DB.Exec(
		query,
		order.ID,
		order.FlightID,
		order.SeatID,
		order.Price,
		order.Status,
	)

	return err
}

func (r *OrderRepo) GetOrderByFlightAndSeat(flightID, seatID string) (*model.Order, error) {

	query := `
	SELECT o.id, o.flight_id, f.flight_number, o.seat_id, s.seat_number, o.price, o.status
	FROM orders o
	JOIN flights f ON o.flight_id = f.id
	JOIN seats s ON o.seat_id = s.id
	WHERE o.flight_id=$1 AND o.seat_id=$2
	`

	var o model.Order

	err := r.DB.QueryRow(query, flightID, seatID).Scan(
		&o.ID,
		&o.FlightID,
		&o.FlightNumber,
		&o.SeatID,
		&o.SeatNumber,
		&o.Price,
		&o.Status,
	)

	if err != nil {
		return nil, err
	}

	return &o, nil
}

func (r *OrderRepo) UpdateOrderStatus(orderID string, status string) error {

	query := `
	UPDATE orders
	SET status = $1
	WHERE id = $2
	`

	_, err := r.DB.Exec(query, status, orderID)
	return err
}
