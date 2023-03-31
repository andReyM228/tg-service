package domain

import "time"

type Car struct {
	ID        int
	Name      string
	Model     string
	Price     int64
	CreatedAt time.Time
}