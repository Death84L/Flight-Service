CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE airplanes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    model TEXT NOT NULL,
    total_seats INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE flights (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    flight_number TEXT NOT NULL,
    source TEXT NOT NULL,
    destination TEXT NOT NULL,
    journey_date DATE NOT NULL,
    departure_time TIMESTAMPTZ NOT NULL,
    arrival_time TIMESTAMPTZ NOT NULL,
    airplane_id UUID NOT NULL REFERENCES airplanes(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE seats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    airplane_id UUID NOT NULL REFERENCES airplanes(id) ON DELETE CASCADE,
    seat_number TEXT NOT NULL,
    price NUMERIC(12,2) NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (airplane_id, seat_number)
);

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    flight_id UUID NOT NULL REFERENCES flights(id) ON DELETE RESTRICT,
    seat_id UUID NOT NULL REFERENCES seats(id) ON DELETE RESTRICT,
    price NUMERIC(12,2) NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (flight_id, seat_id)
);

CREATE INDEX idx_flights_source_destination_date ON flights (source, destination, journey_date);
CREATE INDEX idx_flights_airplane_id ON flights (airplane_id);
CREATE INDEX idx_seats_airplane_id ON seats (airplane_id);
CREATE INDEX idx_orders_flight_id ON orders (flight_id);
CREATE INDEX idx_orders_seat_id ON orders (seat_id);