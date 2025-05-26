package dbrepo

import (
	"context"
	"time"

	"github.com/nelsonmarro/bookings/internal/models"
)

func (m *postgresDbRepo) AllUsers() bool {
	return true
}

func (m *postgresDbRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	newID := 0
	stmt := `
	  INSERT INTO reservations (first_name, last_name, email, phone_number, start_date, end_date, room_id, created_at, updated_at)
		                   VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id
	`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.PhoneNumber,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction inserts a new room restriction into the database.
func (m *postgresDbRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	  insert into room_restrictions (start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at)
		values($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		r.RestrictionID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityByDates checks if there are any room restrictions that overlap with the given start and end dates and room ID.
func (m *postgresDbRepo) SearchAvailabilityByDates(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int
	query := `
		SELECT COUNT(id)
		FROM room_restrictions
		WHERE $1 < end_date AND $2 > start_date
		      and room_id = $3
	`

	err := m.DB.QueryRowContext(ctx, query, start, end, roomID).Scan(&numRows)
	if err != nil {
		return false, err
	}

	return numRows == 0, nil
}
