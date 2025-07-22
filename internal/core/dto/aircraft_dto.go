package dto

import "time"

type Aircraft struct {
    AircraftID           int       `json:"aircraft_id"`
    RegistrationNumber   string    `json:"registration_number"`
    Model                string    `json:"model"`
    Manufacturer         string    `json:"manufacturer"`
    AircraftType         string    `json:"aircraft_type"`
    YearManufactured     int       `json:"year_manufactured"`
    TotalSeats           int       `json:"total_seats"`
    EconomySeats         int       `json:"economy_seats"`
    PremiumEconomySeats  int       `json:"premium_economy_seats"`
    BusinessSeats        int       `json:"business_seats"`
    FirstClassSeats      int       `json:"first_class_seats"`
    MaxRangeKm           int       `json:"max_range_km"`
    CruisingSpeedKmh     int       `json:"cruising_speed_kmh"`
    MaxAltitudeFt        int       `json:"max_altitude_ft"`
    WifiAvailable        bool      `json:"wifi_available"`
    PowerOutletsAvailable bool      `json:"power_outlets_available"`
    EntertainmentSystem  bool      `json:"entertainment_system"`
    AirlineID            int       `json:"airline_id"`
    Status               string    `json:"status"`
    CreatedAt            time.Time `json:"created_at"`
    UpdatedAt            time.Time `json:"updated_at"`
}