-- Create enum type for flight status
CREATE TYPE flight_status AS ENUM ('scheduled', 'delayed', 'cancelled', 'boarding', 'departed', 'arrived', 'diverted');

-- Create enum type for flight class
CREATE TYPE flight_class AS ENUM ('economy', 'premium_economy', 'business', 'first');

-- Create the flights table
CREATE TABLE flights (
    flight_id SERIAL PRIMARY KEY,
    airline_id INT NOT NULL,
    flight_number VARCHAR(10) NOT NULL UNIQUE,
    departure_airport VARCHAR(5) NOT NULL,  -- Using standard IATA airport codes (3-letter) or ICAO codes (4-letter)
    arrival_airport VARCHAR(5) NOT NULL,
    departure_time TIMESTAMP WITH TIME ZONE NOT NULL,
    arrival_time TIMESTAMP WITH TIME ZONE NOT NULL,
    duration_minutes INT NOT NULL,
    duration INTERVAL GENERATED ALWAYS AS (interval '1 minute' * duration_minutes) STORED,  -- Computed column
    stops_count INT NOT NULL DEFAULT 0,
    base_price DECIMAL(10,2) NOT NULL,
    tax_and_fees DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    available_seats INT NOT NULL DEFAULT 0,
    total_seats INT NOT NULL,
     status flight_status NOT NULL DEFAULT 'scheduled',
    class flight_class NOT NULL,
    gate VARCHAR(10) NULL,
    terminal VARCHAR(10) NULL,
    distance INT NULL,  
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_airline FOREIGN KEY (airline_id) REFERENCES airlines(airline_id),
    CONSTRAINT check_times CHECK (departure_time < arrival_time),
    CONSTRAINT check_seats CHECK (available_seats <= total_seats)
);

-- Create index for common queries
CREATE INDEX idx_flights_departure_time ON flights(departure_time);
CREATE INDEX idx_flights_arrival_time ON flights(arrival_time);
CREATE INDEX idx_flights_airports ON flights(departure_airport, arrival_airport);
CREATE INDEX idx_flights_airline ON flights(airline_id);
CREATE INDEX idx_flights_status ON flights(status);


-- Function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger to automatically update the updated_at column
CREATE TRIGGER update_flights_modtime
BEFORE UPDATE ON flights
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();
-- Table airline (fk airline_id)
-- Tạo bảng airlines
CREATE TABLE airlines (
    airline_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    iata_code VARCHAR(2) NOT NULL UNIQUE,   -- 2-character IATA code, e.g., VN
    icao_code VARCHAR(3) NOT NULL UNIQUE,   -- 3-character ICAO code, e.g., HVN
    country VARCHAR(50) NOT NULL,
    logo_url VARCHAR(255) NULL,
    is_low_cost BOOLEAN NOT NULL DEFAULT FALSE,
    description TEXT NULL,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT check_iata_code CHECK (LENGTH(iata_code) = 2),
    CONSTRAINT check_icao_code CHECK (LENGTH(icao_code) = 3)
);

-- Tạo indexes cho các truy vấn phổ biến
CREATE INDEX idx_airlines_country ON airlines(country);
CREATE INDEX idx_airlines_is_low_cost ON airlines(is_low_cost);

-- Function để tự động cập nhật trường updated_at
CREATE OR REPLACE FUNCTION update_airlines_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';
CREATE TRIGGER update_airlines_modtime
BEFORE UPDATE ON airlines
FOR EACH ROW
EXECUTE FUNCTION update_airlines_modified_column();
INSERT INTO airlines (name, iata_code, icao_code, country, is_low_cost, description)
VALUES 
('Vietnam Airlines', 'VN', 'HVN', 'Vietnam', FALSE, 'Flag carrier of Vietnam'),
('VietJet Air', 'VJ', 'VJC', 'Vietnam', TRUE, 'First private airline in Vietnam'),
('Bamboo Airways', 'QH', 'BAV', 'Vietnam', FALSE, 'Vietnamese leisure airline'),
('Singapore Airlines', 'SQ', 'SIA', 'Singapore', FALSE, 'Flag carrier of Singapore'),
('AirAsia', 'AK', 'AXM', 'Malaysia', TRUE, 'Largest low-cost carrier in Asia');
-- Table aircraft (fk airline_id)
CREATE TABLE aircrafts (
    aircraft_id SERIAL PRIMARY KEY,
    airline_id INT NOT NULL,
    model VARCHAR(50) NOT NULL,
    manufacturer VARCHAR(50) NOT NULL,
    registration_number VARCHAR(10) NOT NULL UNIQUE,
    total_seats INT NOT NULL,
    economy_seats INT NOT NULL DEFAULT 0,
    business_seats INT NOT NULL DEFAULT 0,
    first_class_seats INT NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    wifi_available BOOLEAN NOT NULL DEFAULT FALSE,
    inflight_entertainment BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_airline FOREIGN KEY (airline_id) REFERENCES airlines(airline_id),
    CONSTRAINT check_seats CHECK (total_seats = economy_seats + business_seats + first_class_seats),
    CONSTRAINT check_positive_seats CHECK (total_seats > 0 AND economy_seats >= 0 AND business_seats >= 0 AND first_class_seats >= 0),
    CONSTRAINT check_status CHECK (status IN ('active', 'maintenance', 'retired'))
);
CREATE TYPE flight_class AS ENUM ('economy', 'premium_economy', 'business', 'first');

-- Tạo bảng flight_classes
CREATE TABLE flight_classes (
    flight_class_id SERIAL PRIMARY KEY,
    flight_id INT NOT NULL,
    class flight_class NOT NULL,
    base_price DECIMAL(10,2) NOT NULL,
    available_seats INT NOT NULL,
    total_seats INT NOT NULL,
    CONSTRAINT fk_flight FOREIGN KEY (flight_id) REFERENCES flights(flight_id),
    CONSTRAINT check_seats CHECK (available_seats <= total_seats)
);