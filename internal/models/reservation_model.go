package models

import "time"

type ReservationCheck struct {
	StartDate time.Time
	EndDate   time.Time
}
