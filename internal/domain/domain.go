package domain

import "time"

type Car struct {
	ID        int
	Name      string
	Model     string
	Price     int64
	Image     string
	CreatedAt time.Time
}

type Cars struct {
	Cars []Car
}

type User struct {
	ID        int
	Name      string
	Surname   string
	Phone     string
	Email     string
	Password  string
	ChatID    int64
	Cars      []Car
	CreatedAt time.Time
}

type CarCharacteristics struct {
	Engine       string `json:"engine,omitempty"`
	Power        string `json:"power,omitempty"`
	Acceleration string `json:"acceleration,omitempty"`
	TopSpeed     string `json:"top_speed,omitempty"`
	FuelEconomy  string `json:"fuel_economy,omitempty"`
	Transmission string `json:"transmission,omitempty"`
	DriveType    string `json:"drive_type,omitempty"`
}
