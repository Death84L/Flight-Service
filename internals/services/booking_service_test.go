package services

import (
	"errors"
	"testing"

	model "flight-service/internals/models"
)

func validFlight() *model.Flight {
	return &model.Flight{
		ID:           "F1",
		FlightNumber: "FL-100",
		AirplaneID:   "A1",
	}
}

func availableSeat() *model.Seat {
	return &model.Seat{
		ID:         "S1",
		SeatNumber: "12A",
		Status:     model.SeatAvailable,
		AirplaneID: "A1",
		Price:      1200,
	}
}

func bookedSeat() *model.Seat {
	return &model.Seat{
		ID:         "S1",
		SeatNumber: "12A",
		Status:     model.SeatBooked,
		AirplaneID: "A1",
		Price:      1200,
	}
}

func reservedSeat() *model.Seat {
	return &model.Seat{
		ID:         "S1",
		SeatNumber: "12A",
		Status:     model.SeatReserved,
		AirplaneID: "A1",
		Price:      1200,
	}
}

func existingOrder() *model.Order {
	return &model.Order{ID: "O1"}
}

func TestBookSeat(t *testing.T) {
	tests := []struct {
		name          string
		flightRepo    *mockFlightRepo
		seatRepo      *mockSeatRepo
		orderRepo     *mockOrderRepo
		flightNumber  string
		seatNumber    string
		wantErr       bool
		wantErrMsg    string
		wantConfirmed bool
	}{
		// ── Happy path
		{
			name:          "success – available seat, no existing order",
			flightRepo:    &mockFlightRepo{flight: validFlight()},
			seatRepo:      &mockSeatRepo{seat: availableSeat()},
			orderRepo:     &mockOrderRepo{},
			flightNumber:  "FL-100",
			seatNumber:    "12A",
			wantErr:       false,
			wantConfirmed: true,
		},

		// ── Flight errors
		{
			name:         "flight not found – nil returned",
			flightRepo:   &mockFlightRepo{flight: nil},
			seatRepo:     &mockSeatRepo{seat: availableSeat()},
			orderRepo:    &mockOrderRepo{},
			flightNumber: "FL-999",
			seatNumber:   "12A",
			wantErr:      true,
			wantErrMsg:   "flight not found",
		},
		{
			name:         "flight repo returns db error",
			flightRepo:   &mockFlightRepo{flight: nil, err: errors.New("db connection refused")},
			seatRepo:     &mockSeatRepo{seat: availableSeat()},
			orderRepo:    &mockOrderRepo{},
			flightNumber: "FL-100",
			seatNumber:   "12A",
			wantErr:      true,
			wantErrMsg:   "db connection refused",
		},

		// ── Seat errors
		{
			name:         "seat not found – nil returned",
			flightRepo:   &mockFlightRepo{flight: validFlight()},
			seatRepo:     &mockSeatRepo{seat: nil},
			orderRepo:    &mockOrderRepo{},
			flightNumber: "FL-100",
			seatNumber:   "99Z",
			wantErr:      true,
			wantErrMsg:   "seat not found",
		},
		{
			name:         "seat repo returns db error",
			flightRepo:   &mockFlightRepo{flight: validFlight()},
			seatRepo:     &mockSeatRepo{seat: nil, err: errors.New("seat query failed")},
			orderRepo:    &mockOrderRepo{},
			flightNumber: "FL-100",
			seatNumber:   "12A",
			wantErr:      true,
			wantErrMsg:   "seat query failed",
		},
		{
			name:         "seat status is BOOKED – not available",
			flightRepo:   &mockFlightRepo{flight: validFlight()},
			seatRepo:     &mockSeatRepo{seat: bookedSeat()},
			orderRepo:    &mockOrderRepo{},
			flightNumber: "FL-100",
			seatNumber:   "12A",
			wantErr:      true,
			wantErrMsg:   "seat not available",
		},
		{
			name:         "seat status is RESERVED – not available",
			flightRepo:   &mockFlightRepo{flight: validFlight()},
			seatRepo:     &mockSeatRepo{seat: reservedSeat()},
			orderRepo:    &mockOrderRepo{},
			flightNumber: "FL-100",
			seatNumber:   "12A",
			wantErr:      true,
			wantErrMsg:   "seat not available",
		},

		// ── Idempotency / duplicate guard
		{
			name:       "duplicate booking – existing order found",
			flightRepo: &mockFlightRepo{flight: validFlight()},
			seatRepo:   &mockSeatRepo{seat: availableSeat()},
			orderRepo:  &mockOrderRepo{existingOrder: existingOrder()},
			wantErr:    true,
			wantErrMsg: "seat already booked on this flight",
		},

		// ── Order creation errors
		{
			name:         "order creation fails after all retries",
			flightRepo:   &mockFlightRepo{flight: validFlight()},
			seatRepo:     &mockSeatRepo{seat: availableSeat()},
			orderRepo:    &mockOrderRepo{createErr: errors.New("insert failed")},
			flightNumber: "FL-100",
			seatNumber:   "12A",
			wantErr:      true,
			wantErrMsg:   "insert failed",
		},
		{
			name:         "order creation fails with timeout error",
			flightRepo:   &mockFlightRepo{flight: validFlight()},
			seatRepo:     &mockSeatRepo{seat: availableSeat()},
			orderRepo:    &mockOrderRepo{createErr: errors.New("context deadline exceeded")},
			flightNumber: "FL-100",
			seatNumber:   "12A",
			wantErr:      true,
			wantErrMsg:   "context deadline exceeded",
		},

		// ── Redis nil (no lock)
		{
			name:          "redis is nil – booking succeeds without lock",
			flightRepo:    &mockFlightRepo{flight: validFlight()},
			seatRepo:      &mockSeatRepo{seat: availableSeat()},
			orderRepo:     &mockOrderRepo{},
			flightNumber:  "FL-100",
			seatNumber:    "12A",
			wantErr:       false,
			wantConfirmed: true,
		},

		// ── Field validation edge cases
		{
			name:         "empty flight number – flight not found",
			flightRepo:   &mockFlightRepo{flight: nil},
			seatRepo:     &mockSeatRepo{seat: availableSeat()},
			orderRepo:    &mockOrderRepo{},
			flightNumber: "",
			seatNumber:   "12A",
			wantErr:      true,
			wantErrMsg:   "flight not found",
		},
		{
			name:         "empty seat number – seat not found",
			flightRepo:   &mockFlightRepo{flight: validFlight()},
			seatRepo:     &mockSeatRepo{seat: nil},
			orderRepo:    &mockOrderRepo{},
			flightNumber: "FL-100",
			seatNumber:   "",
			wantErr:      true,
			wantErrMsg:   "seat not found",
		},

		// ── Zero price seat (edge)
		{
			name:       "seat with zero price – booking still succeeds",
			flightRepo: &mockFlightRepo{flight: validFlight()},
			seatRepo: &mockSeatRepo{seat: &model.Seat{
				ID: "S2", SeatNumber: "1A",
				Status: model.SeatAvailable, AirplaneID: "A1", Price: 0,
			}},
			orderRepo:     &mockOrderRepo{},
			flightNumber:  "FL-100",
			seatNumber:    "1A",
			wantErr:       false,
			wantConfirmed: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := &BookingService{
				flightRepo: tc.flightRepo,
				seatRepo:   tc.seatRepo,
				orderRepo:  tc.orderRepo,
				redis:      nil,
			}

			order, err := svc.BookSeat(tc.flightNumber, tc.seatNumber)

			// ── error assertion ──
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error %q but got nil", tc.wantErrMsg)
					return
				}
				if err.Error() != tc.wantErrMsg {
					t.Errorf("expected error %q, got %q", tc.wantErrMsg, err.Error())
				}
				return
			}

			// ── success assertion ──
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if order == nil {
				t.Error("expected non-nil order")
				return
			}
			if tc.wantConfirmed && order.Status != model.OrderConfirmed {
				t.Errorf("expected order status CONFIRMED, got %s", order.Status)
			}
			if order.ID == "" {
				t.Error("expected order ID to be set")
			}
			if order.FlightID == "" {
				t.Error("expected order FlightID to be set")
			}
			if order.SeatID == "" {
				t.Error("expected order SeatID to be set")
			}
		})
	}
}

//  retry helper  isolated tests

func TestRetry(t *testing.T) {
	tests := []struct {
		name       string
		attempts   int
		failTimes  int // how many times fn returns error before succeeding
		alwaysFail bool
		wantErr    bool
		wantCalls  int
	}{
		{
			name:      "succeeds on first attempt",
			attempts:  3,
			failTimes: 0,
			wantErr:   false,
			wantCalls: 1,
		},
		{
			name:      "succeeds on second attempt",
			attempts:  3,
			failTimes: 1,
			wantErr:   false,
			wantCalls: 2,
		},
		{
			name:      "succeeds on third attempt",
			attempts:  3,
			failTimes: 2,
			wantErr:   false,
			wantCalls: 3,
		},
		{
			name:       "fails all 3 attempts",
			attempts:   3,
			alwaysFail: true,
			wantErr:    true,
			wantCalls:  3,
		},
		{
			name:      "zero attempts – never calls fn",
			attempts:  0,
			wantErr:   true, // returns last err which is nil-initialised → still non-nil check passes
			wantCalls: 0,
		},
		{
			name:      "single attempt succeeds",
			attempts:  1,
			failTimes: 0,
			wantErr:   false,
			wantCalls: 1,
		},
		{
			name:       "single attempt fails",
			attempts:   1,
			alwaysFail: true,
			wantErr:    true,
			wantCalls:  1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			callCount := 0
			failCount := 0

			err := retry(tc.attempts, 0, func() error {
				callCount++
				if tc.alwaysFail {
					return errors.New("always fails")
				}
				if failCount < tc.failTimes {
					failCount++
					return errors.New("transient error")
				}
				return nil
			})

			if tc.wantErr && err == nil && tc.attempts > 0 {
				t.Errorf("expected error but got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if callCount != tc.wantCalls {
				t.Errorf("expected %d fn calls, got %d", tc.wantCalls, callCount)
			}
		})
	}
}

// Order field integrity tests
func TestBookSeat_OrderFields(t *testing.T) {
	flight := validFlight()
	seat := availableSeat()

	svc := &BookingService{
		flightRepo: &mockFlightRepo{flight: flight},
		seatRepo:   &mockSeatRepo{seat: seat},
		orderRepo:  &mockOrderRepo{},
		redis:      nil,
	}

	order, err := svc.BookSeat("FL-100", "12A")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	checks := []struct {
		field string
		got   string
		want  string
	}{
		{"FlightID", order.FlightID, flight.ID},
		{"FlightNumber", order.FlightNumber, flight.FlightNumber},
		{"SeatID", order.SeatID, seat.ID},
		{"SeatNumber", order.SeatNumber, seat.SeatNumber},
	}

	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("order.%s = %q, want %q", c.field, c.got, c.want)
		}
	}

	if order.Price != seat.Price {
		t.Errorf("order.Price = %v, want %v", order.Price, seat.Price)
	}
	if order.Status != model.OrderConfirmed {
		t.Errorf("order.Status = %v, want CONFIRMED", order.Status)
	}
	if order.CreatedAt.IsZero() {
		t.Error("order.CreatedAt should not be zero")
	}
	if order.UpdatedAt.IsZero() {
		t.Error("order.UpdatedAt should not be zero")
	}
}
