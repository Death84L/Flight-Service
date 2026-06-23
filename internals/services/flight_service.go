package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	model "flight-service/internals/models"
	"flight-service/internals/redis"
	"flight-service/internals/repository/interfaces"
)

// FlightService handles ONLY flight search (READ operations)
type FlightService struct {
	flightRepo interfaces.FlightRepository
	redis      *redis.RedisClient
}

func NewFlightService(
	flightRepo interfaces.FlightRepository,
	redisClient *redis.RedisClient,
) *FlightService {
	return &FlightService{
		flightRepo: flightRepo,
		redis:      redisClient,
	}
}

const (
	defaultFlightLimit  = 10
	defaultFlightOffset = 0
	cacheTTL            = 10 * time.Minute
)

func (s *FlightService) SearchFlights(
	ctx context.Context,
	source string,
	destination string,
	date string,
	limit int,
	offset int,
) ([]model.Flight, error) {

	// -----------------------------
	// timeout protection
	// -----------------------------
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	log.Printf("[SEARCH_FLIGHTS] source=%s destination=%s date=%s", source, destination, date)

	if limit <= 0 {
		limit = defaultFlightLimit
	}
	if offset < 0 {
		offset = defaultFlightOffset
	}

	cacheKey := fmt.Sprintf(
		"flights:%s:%s",
		source,
		destination,
	)

	// cache first
	if s.redis != nil {
		cached, err := s.redis.Client.Get(ctx, cacheKey).Result()
		if err == nil && cached != "" {
			var flights []model.Flight
			if err := json.Unmarshal([]byte(cached), &flights); err == nil {
				log.Printf("[CACHE_HIT] %s", cacheKey)
				return flights, nil
			}
		}
	}

	log.Printf("[CACHE_MISS] %s", cacheKey)

	//Fetch from db
	flights, err := s.flightRepo.GetFlights(ctx, source, destination, date, limit, offset)
	if err != nil {
		log.Printf("[SEARCH_FLIGHTS_ERROR] %v", err)
		return nil, fmt.Errorf("search flights failed: %w", err)
	}

	if flights == nil {
		return []model.Flight{}, nil
	}

	//store in cache
	if s.redis != nil {
		data, err := json.Marshal(flights)
		if err == nil {
			_ = s.redis.Client.Set(ctx, cacheKey, data, cacheTTL).Err()
		}
	}

	return flights, nil
}
