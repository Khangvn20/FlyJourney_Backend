package dto

import "time"

type Airline struct {
    AirlineID  int       `json:"airline_id"`
    Name       string    `json:"name"`
    IATACode   string    `json:"iata_code"`
    ICAOCode   string    `json:"icao_code"`
    Country    string    `json:"country"`
    LogoURL    string    `json:"logo_url"`
    IsLowCost  bool      `json:"is_low_cost"`
    Description string    `json:"description"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}