package postgres

import (
	"database/sql"
	model "flight-service/internals/models"
)

type SeatRepo struct {
	DB *sql.DB
}

func NewSeatRepo(db *sql.DB) *SeatRepo {
	return &SeatRepo{DB: db}
}

func (r *SeatRepo) GetSeatsByAirplaneID(airplaneID string) ([]model.Seat, error) {

	query := `
	SELECT id, airplane_id, seat_number, price, status
	FROM seats
	WHERE airplane_id=$1
	`

	rows, err := r.DB.Query(query, airplaneID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []model.Seat

	for rows.Next() {
		var s model.Seat
		err := rows.Scan(
			&s.ID,
			&s.AirplaneID,
			&s.SeatNumber,
			&s.Price,
			&s.Status,
		)
		if err != nil {
			return nil, err
		}

		seats = append(seats, s)
	}

	return seats, nil
}

func (r *SeatRepo) GetSeatByAirplaneIDAndSeatNumber(airplaneID, seatNumber string) (*model.Seat, error) {

	query := `
	SELECT id, airplane_id, seat_number, price, status
	FROM seats
	WHERE airplane_id=$1 AND seat_number=$2
	`

	var s model.Seat
	err := r.DB.QueryRow(query, airplaneID, seatNumber).Scan(
		&s.ID,
		&s.AirplaneID,
		&s.SeatNumber,
		&s.Price,
		&s.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &s, nil
}

func (r *SeatRepo) UpdateSeatStatus(seatID string, status string) error {

	query := `
	UPDATE seats
	SET status=$1
	WHERE id=$2
	`

	_, err := r.DB.Exec(query, status, seatID)
	return err
}
