package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	model "flight-service/internals/models"
	"flight-service/internals/services"
)

type Controller struct {
	flightService  *services.FlightService
	bookingService *services.BookingService
}

func NewController(f *services.FlightService, b *services.BookingService) *Controller {
	return &Controller{
		flightService:  f,
		bookingService: b,
	}
}

func (c *Controller) SearchFlights(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, NewBadRequest("method not allowed"))
		return
	}
	source := r.URL.Query().Get("source")
	destination := r.URL.Query().Get("destination")
	date := r.URL.Query().Get("date")

	if source == "" || destination == "" {
		WriteError(w, NewBadRequest("source and destination required"))
		return
	}

	limit := model.DefaultFlightLimit
	offset := model.DefaultFlightOffset

	if q := r.URL.Query().Get("limit"); q != "" {
		if v, err := strconv.Atoi(q); err == nil {
			limit = v
		}
	}
	if q := r.URL.Query().Get("offset"); q != "" {
		if v, err := strconv.Atoi(q); err == nil {
			offset = v
		}
	}

	flights, err := c.flightService.SearchFlights(r.Context(), source, destination, date, limit, offset)
	if err != nil {
		WriteError(w, NewInternalError(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flights)
}

func (c *Controller) BookSeat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, NewBadRequest("method not allowed"))
		return
	}
	var req struct {
		FlightNumber string `json:"flight_number"`
		SeatNumber   string `json:"seat_number"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		WriteError(w, NewBadRequest("invalid request body"))
		return
	}

	if req.FlightNumber == "" || req.SeatNumber == "" {
		WriteError(w, NewBadRequest("missing required fields"))
		return
	}

	order, err := c.bookingService.BookSeat(
		req.FlightNumber,
		req.SeatNumber,
	)

	if err != nil {
		WriteError(w, NewInternalError(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
