package routes

import (
	"net/http"

	"flight-service/internals/handlers"
)

func RegisterRoutes(controller *handlers.Controller) {
	// API VERSION PREFIX
	v1 := "/api/v1"

	// GET - flight search
	http.HandleFunc(v1+"/flights", func(w http.ResponseWriter, r *http.Request) {
		controller.SearchFlights(w, r)
	})

	// POST - book seat
	http.HandleFunc(v1+"/book", func(w http.ResponseWriter, r *http.Request) {
		controller.BookSeat(w, r)
	})
}
