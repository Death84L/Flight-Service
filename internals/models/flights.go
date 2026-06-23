package model

import "time"

type Flight struct {
	ID string `json:"id"`

	FlightNumber string `json:"flight_number"`

	Source      string `json:"source"`
	Destination string `json:"destination"`

	JourneyDate time.Time `json:"journey_date"`

	DepartureTime time.Time `json:"departure_time"`
	ArrivalTime   time.Time `json:"arrival_time"`

	AirplaneID string `json:"airplane_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
