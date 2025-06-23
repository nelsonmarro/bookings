package dbrepo

import (
	"errors"
	"time"

	"github.com/nelsonmarro/bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	if res.RoomID == 2 {
		return 0, errors.New("room not available")
	}
	return 1, nil
}

// InsertRoomRestriction inserts a new room restriction into the database.
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	return nil
}

// SearchAvailabilityByDatesByRoomID searches for room availability based on the provided start and end dates, and room ID.
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	return false, nil
}

// SearchAvailabilityForAllRooms checks the availability of all rooms for the given date range.
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// GetRoomByID retrieves a room by its ID from the database.
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room

	if id > 2 {
		return room, errors.New("room not found")
	}

	return room, nil
}

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	return models.User{}, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

// Authenticate checks if the user exists and if the password is correct.
func (m *testDBRepo) Authenticate(email, password string) (int, string, error) {
	return 0, "", errors.New("not implemented")
}

func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}

func (m *testDBRepo) NewReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}

func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	return models.Reservation{}, nil
}

func (m *testDBRepo) UpdateReservation(res models.Reservation) error {
	return nil
}

func (m *testDBRepo) DeleteReservation(id int) error {
	return nil
}

func (m *testDBRepo) UpdateProcessedForReservation(id int, processed bool) error {
	return nil
}

func (m *testDBRepo) AllRooms() ([]models.Room, error) {
	var rooms []models.Room

	return rooms, nil
}

func (m *testDBRepo) GetRestrictionsForRoomByDates(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {
	return nil, nil
}

func (m *testDBRepo) InsertBlockForRoom(id int, start time.Time) error {
	return nil
}

func (m *testDBRepo) DeleteBlockByID(id int) error {
	return nil
}
