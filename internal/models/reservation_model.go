package models

import "time"

// Reservation is the reservation model
type Reservation struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
	StartDate   time.Time
	EndDate     time.Time
	RoomID      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Processed   bool
	Room        Room
}
