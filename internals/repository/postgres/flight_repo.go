package postgres

import (
	"context"
	"database/sql"
	"fmt"

	model "flight-service/internals/models"
)

type FlightRepo struct {
	DB *sql.DB
}

func NewFlightRepo(db *sql.DB) *FlightRepo {
	return &FlightRepo{DB: db}
}

// ===============================
// GET FLIGHTS (SEARCH)
// ===============================
func (r *FlightRepo) GetFlights(
	ctx context.Context,
	source, destination string,
	date string,
	limit, offset int,
) ([]model.Flight, error) {

	query := `
	SELECT id, flight_number, source, destination, journey_date,
	       departure_time, arrival_time, airplane_id
	FROM flights
	WHERE source = $1 AND destination = $2
	`

	args := []interface{}{source, destination}
	argPos := 3

	// optional filter (null-safe)
	if date != "" {
		query += fmt.Sprintf(" AND journey_date = $%d", argPos)
		args = append(args, date)
		argPos++
	}

	query += fmt.Sprintf(
		" ORDER BY departure_time LIMIT $%d OFFSET $%d",
		argPos, argPos+1,
	)

	args = append(args, limit, offset)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("GetFlights query failed: %w", err)
	}
	defer rows.Close()

	var flights []model.Flight

	for rows.Next() {
		var f model.Flight

		err := rows.Scan(
			&f.ID,
			&f.FlightNumber,
			&f.Source,
			&f.Destination,
			&f.JourneyDate,
			&f.DepartureTime,
			&f.ArrivalTime,
			&f.AirplaneID,
		)

		if err != nil {
			return nil, fmt.Errorf("GetFlights scan failed: %w", err)
		}

		flights = append(flights, f)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetFlights rows iteration failed: %w", err)
	}

	return flights, nil
}

// GET FLIGHT BY NUMBER
// ===============================
func (r *FlightRepo) GetFlightByNumber(
	ctx context.Context,
	flightNumber string,
) (*model.Flight, error) {

	query := `
	SELECT id, flight_number, source, destination, journey_date,
	       departure_time, arrival_time, airplane_id
	FROM flights
	WHERE flight_number = $1
	`

	var f model.Flight

	err := r.DB.QueryRowContext(ctx, query, flightNumber).Scan(
		&f.ID,
		&f.FlightNumber,
		&f.Source,
		&f.Destination,
		&f.JourneyDate,
		&f.DepartureTime,
		&f.ArrivalTime,
		&f.AirplaneID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("GetFlightByNumber failed: %w", err)
	}

	return &f, nil
}
