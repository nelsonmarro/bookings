package repository

import "github.com/nelsonmarro/bookings/internal/models"

type DataBaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
}
