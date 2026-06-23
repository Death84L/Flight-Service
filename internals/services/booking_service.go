package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"

	model "flight-service/internals/models"
	"flight-service/internals/redis"
	"flight-service/internals/repository/interfaces"
)

// BookingService handles ONLY booking flow (WRITE operations)
type BookingService struct {
	flightRepo interfaces.FlightRepository
	seatRepo   interfaces.SeatRepository
	orderRepo  interfaces.OrderRepository
	redis      *redis.RedisClient
}

// Constructor
func NewBookingService(
	flightRepo interfaces.FlightRepository,
	seatRepo interfaces.SeatRepository,
	orderRepo interfaces.OrderRepository,
	redisClient *redis.RedisClient,
) *BookingService {

	return &BookingService{
		flightRepo: flightRepo,
		seatRepo:   seatRepo,
		orderRepo:  orderRepo,
		redis:      redisClient,
	}
}

// BOOK SEAT
func (s *BookingService) BookSeat(
	flightNumber string,
	seatNumber string,
) (*model.Order, error) {

	log.Printf("[BOOKING_START] flight=%s seat=%s", flightNumber, seatNumber)

	flight, err := s.flightRepo.GetFlightByNumber(context.Background(), flightNumber)
	if err != nil {
		log.Printf("[ERROR] flight fetch failed: %v", err)
		return nil, err
	}
	if flight == nil {
		log.Printf("[ERROR] flight not found")
		return nil, errors.New("flight not found")
	}

	selectedSeat, err := s.seatRepo.GetSeatByAirplaneIDAndSeatNumber(flight.AirplaneID, seatNumber)
	if err != nil {
		log.Printf("[ERROR] seat fetch failed: %v", err)
		return nil, err
	}
	if selectedSeat == nil {
		log.Printf("[ERROR] seat not found")
		return nil, errors.New("seat not found")
	}

	seatID := selectedSeat.ID

	existingOrder, _ := s.orderRepo.GetOrderByFlightAndSeat(flight.ID, seatID)
	if existingOrder != nil {
		log.Printf("[ERROR] seat already booked")
		return nil, errors.New("seat already booked on this flight")
	}

	var lockKey string
	if s.redis != nil {
		lockKey = "seat_lock:" + flight.AirplaneID + ":" + seatID

		ok, err := s.redis.Client.SetNX(
			s.redis.Ctx,
			lockKey,
			"locked",
			5*time.Minute,
		).Result()

		if err != nil {
			log.Printf("[ERROR] redis lock error: %v", err)
			return nil, err
		}

		if !ok {
			log.Printf("[ERROR] seat locked by another request")
			return nil, errors.New("seat already locked")
		}

		defer func() {
			_ = s.redis.Client.Del(s.redis.Ctx, lockKey)
		}()
	}

	if selectedSeat.Status != model.SeatAvailable {
		log.Printf("[ERROR] seat not available")
		return nil, errors.New("seat not available")
	}

	order := model.Order{
		ID:           uuid.New().String(),
		FlightID:     flight.ID,
		FlightNumber: flight.FlightNumber,
		SeatID:       seatID,
		SeatNumber:   selectedSeat.SeatNumber,
		Price:        selectedSeat.Price,
		Status:       model.OrderPending,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// RETRY WRAPPED CREATE ORDER
	err = retry(3, 100*time.Millisecond, func() error {
		return s.orderRepo.CreateOrder(order)
	})

	if err != nil {
		log.Printf("[ERROR] order creation failed after retry: %v", err)
		return nil, err
	}

	// RETRY WRAPPED SEAT UPDATE
	err = retry(3, 100*time.Millisecond, func() error {
		return s.seatRepo.UpdateSeatStatus(seatID, string(model.SeatBooked))
	})

	if err != nil {
		log.Printf("[ERROR] seat update failed after retry: %v", err)
		return nil, err
	}

	if s.redis != nil {
		_ = s.redis.Client.Del(s.redis.Ctx, lockKey)
	}

	order.Status = model.OrderConfirmed

	log.Printf("[BOOKING_SUCCESS] order_id=%s", order.ID)

	return &order, nil
}

// RETRY HELPER
func retry(attempts int, delay time.Duration, fn func() error) error {
	var err error

	for i := 0; i < attempts; i++ {
		err = fn()
		if err == nil {
			return nil
		}

		log.Printf("[RETRY] attempt=%d error=%v", i+1, err)
		time.Sleep(delay)
	}

	return err
}
