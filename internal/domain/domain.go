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
