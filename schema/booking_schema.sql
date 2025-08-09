CREATE TABLE bookings (
    booking_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    flight_id INTEGER REFERENCES flights(flight_id),
    booking_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE booking_details (
    booking_detail_id SERIAL PRIMARY KEY,
    booking_id INTEGER REFERENCES bookings(booking_id),
    passenger_name VARCHAR(100) NOT NULL,
    passenger_age INTEGER,
    passenger_gender VARCHAR(10),
    flight_class_id INTEGER REFERENCES flight_classes(flight_class_id),
    seat_id INTEGER REFERENCES seats(seat_id),
    price DECIMAL(10,2) NOT NULL
);

CREATE TABLE payments (
    payment_id SERIAL PRIMARY KEY,
    booking_id INTEGER REFERENCES bookings(booking_id),
    amount DECIMAL(10,2) NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    transaction_id VARCHAR(100),
    paid_at TIMESTAMP
);