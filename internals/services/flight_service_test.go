package services

import (
	"context"
	"errors"
	"testing"

	model "flight-service/internals/models"
)

// SearchFlights – table-driven tests
func TestSearchFlights(t *testing.T) {
	tests := []struct {
		name        string
		flightRepo  *mockFlightRepo
		source      string
		destination string
		date        string
		limit       int
		offset      int
		wantLen     int
		wantErr     bool
		wantErrMsg  string
	}{
		// ── Happy path
		{
			name: "success – single flight returned",
			flightRepo: &mockFlightRepo{
				flights: []model.Flight{
					{ID: "1", Source: "BLR", Destination: "DEL"},
				},
			},
			source: "BLR", destination: "DEL", date: "2026-01-01",
			limit: 10, offset: 0,
			wantLen: 1,
		},
		{
			name: "success – multiple flights returned",
			flightRepo: &mockFlightRepo{
				flights: []model.Flight{
					{ID: "1", Source: "BLR", Destination: "DEL"},
					{ID: "2", Source: "BLR", Destination: "DEL"},
					{ID: "3", Source: "BLR", Destination: "DEL"},
				},
			},
			source: "BLR", destination: "DEL", date: "",
			limit: 10, offset: 0,
			wantLen: 3,
		},
		{
			name:       "success – empty result from repo",
			flightRepo: &mockFlightRepo{flights: []model.Flight{}},
			source:     "BLR", destination: "DEL", date: "",
			limit: 10, offset: 0,
			wantLen: 0,
		},
		{
			name:       "success – nil result from repo returns empty slice",
			flightRepo: &mockFlightRepo{flights: nil},
			source:     "BLR", destination: "DEL", date: "",
			limit: 10, offset: 0,
			wantLen: 0,
		},

		// ── Repo errors
		{
			name:       "repo returns db error",
			flightRepo: &mockFlightRepo{err: errors.New("db error")},
			source:     "BLR", destination: "DEL", date: "",
			limit: 10, offset: 0,
			wantErr:    true,
			wantErrMsg: "search flights failed: db error",
		},
		{
			name:       "repo returns connection refused",
			flightRepo: &mockFlightRepo{err: errors.New("connection refused")},
			source:     "BLR", destination: "DEL", date: "",
			limit: 10, offset: 0,
			wantErr:    true,
			wantErrMsg: "search flights failed: connection refused",
		},
		{
			name:       "repo returns context deadline exceeded",
			flightRepo: &mockFlightRepo{err: errors.New("context deadline exceeded")},
			source:     "BLR", destination: "DEL", date: "",
			limit: 10, offset: 0,
			wantErr:    true,
			wantErrMsg: "search flights failed: context deadline exceeded",
		},

		// ── Limit / offset edge cases
		{
			name: "limit zero – defaults to 10",
			flightRepo: &mockFlightRepo{
				flights: []model.Flight{{ID: "1"}},
			},
			source: "BLR", destination: "DEL", date: "",
			limit: 0, offset: 0, // should default to 10 internally
			wantLen: 1,
		},
		{
			name: "negative limit – defaults to 10",
			flightRepo: &mockFlightRepo{
				flights: []model.Flight{{ID: "1"}},
			},
			source: "BLR", destination: "DEL", date: "",
			limit: -5, offset: 0,
			wantLen: 1,
		},
		{
			name: "negative offset – defaults to 0",
			flightRepo: &mockFlightRepo{
				flights: []model.Flight{{ID: "1"}},
			},
			source: "BLR", destination: "DEL", date: "",
			limit: 10, offset: -1,
			wantLen: 1,
		},
		{
			name: "large limit value",
			flightRepo: &mockFlightRepo{
				flights: []model.Flight{{ID: "1"}, {ID: "2"}},
			},
			source: "BLR", destination: "DEL", date: "",
			limit: 1000, offset: 0,
			wantLen: 2,
		},
		{
			name: "non-zero offset",
			flightRepo: &mockFlightRepo{
				flights: []model.Flight{{ID: "2"}},
			},
			source: "BLR", destination: "DEL", date: "",
			limit: 10, offset: 1,
			wantLen: 1,
		},

		// ── Date filter edge cases
		{
			name: "with date filter",
			flightRepo: &mockFlightRepo{
				flights: []model.Flight{{ID: "1"}},
			},
			source: "BLR", destination: "DEL", date: "2026-06-01",
			limit: 10, offset: 0,
			wantLen: 1,
		},
		{
			name:       "with date filter – no matching flights",
			flightRepo: &mockFlightRepo{flights: []model.Flight{}},
			source:     "BLR", destination: "DEL", date: "2099-01-01",
			limit: 10, offset: 0,
			wantLen: 0,
		},
		{
			name:       "empty date – no filter applied",
			flightRepo: &mockFlightRepo{flights: []model.Flight{{ID: "1"}}},
			source:     "BLR", destination: "DEL", date: "",
			limit: 10, offset: 0,
			wantLen: 1,
		},

		// ── Source / destination edge cases
		{
			name:       "same source and destination",
			flightRepo: &mockFlightRepo{flights: []model.Flight{}},
			source:     "BLR", destination: "BLR", date: "",
			limit: 10, offset: 0,
			wantLen: 0,
		},
		{
			name:       "empty source and destination",
			flightRepo: &mockFlightRepo{flights: []model.Flight{}},
			source:     "", destination: "", date: "",
			limit: 10, offset: 0,
			wantLen: 0,
		},

		// ── Pagination result sets
		{
			name: "pagination – first page of 2",
			flightRepo: &mockFlightRepo{
				flights: []model.Flight{{ID: "1"}, {ID: "2"}},
			},
			source: "BLR", destination: "DEL", date: "",
			limit: 2, offset: 0,
			wantLen: 2,
		},
		{
			name: "pagination – second page returns 1",
			flightRepo: &mockFlightRepo{
				flights: []model.Flight{{ID: "3"}},
			},
			source: "BLR", destination: "DEL", date: "",
			limit: 2, offset: 2,
			wantLen: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewFlightService(tc.flightRepo, nil)

			res, err := svc.SearchFlights(
				context.Background(),
				tc.source,
				tc.destination,
				tc.date,
				tc.limit,
				tc.offset,
			)

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
			if len(res) != tc.wantLen {
				t.Errorf("expected %d flights, got %d", tc.wantLen, len(res))
			}
		})
	}
}

// Flight field integrity
func TestSearchFlights_FieldIntegrity(t *testing.T) {
	expected := model.Flight{
		ID:           "F1",
		FlightNumber: "FL-100",
		Source:       "BLR",
		Destination:  "DEL",
		AirplaneID:   "A1",
	}

	svc := NewFlightService(&mockFlightRepo{
		flights: []model.Flight{expected},
	}, nil)

	res, err := svc.SearchFlights(context.Background(), "BLR", "DEL", "", 10, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("expected 1 flight, got %d", len(res))
	}

	got := res[0]
	checks := []struct {
		field string
		got   string
		want  string
	}{
		{"ID", got.ID, expected.ID},
		{"FlightNumber", got.FlightNumber, expected.FlightNumber},
		{"Source", got.Source, expected.Source},
		{"Destination", got.Destination, expected.Destination},
		{"AirplaneID", got.AirplaneID, expected.AirplaneID},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("flight.%s = %q, want %q", c.field, c.got, c.want)
		}
	}
}
