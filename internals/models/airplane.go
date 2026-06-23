package model

import "time"

type Airplane struct {
	ID string `json:"id"`

	Model string `json:"model"`

	TotalSeats int `json:"total_seats"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
