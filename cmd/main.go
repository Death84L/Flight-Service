package main

import (
	"fmt"
	"log"
	"net/http"

	"flight-service/internals/database"
	"flight-service/internals/handlers"
	"flight-service/internals/redis"
	"flight-service/internals/repository/postgres"
	"flight-service/internals/routes"
	"flight-service/internals/services"
)

func main() {

	fmt.Println("Starting the server...")

	// 1. PostgreSQL connection
	pg := database.NewPostgres(
		"HOST",
		"PORT",
		"USER",
		"PWD",
		"DB_NAME",
	)

	// 2. Redis connection
	redisClient := redis.NewRedis("localhost", "6379")

	// 3. Initialize repositories
	flightRepo := postgres.NewFlightRepo(pg.DB)
	seatRepo := postgres.NewSeatRepo(pg.DB)
	orderRepo := postgres.NewOrderRepo(pg.DB)

	// 4. Initialize services
	flightService := services.NewFlightService(flightRepo, redisClient)
	bookingService := services.NewBookingService(flightRepo, seatRepo, orderRepo, redisClient)

	// 5. Initialize controller
	controller := handlers.NewController(flightService, bookingService)

	// 6. Register routes
	routes.RegisterRoutes(controller)

	// 7. Start server
	port := ":8080"
	fmt.Printf("Server is running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
