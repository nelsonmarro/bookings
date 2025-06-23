package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/nelsonmarro/bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
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

// SearchAvailabilityByDatesByRoomID searches for room availability based on the provided start and end dates, and room ID.
func (m *postgresDbRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
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

// SearchAvailabilityForAllRooms checks the availability of all rooms for the given date range.
func (m *postgresDbRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `
select
	r.id, r.room_name
from
rooms r
where
r.id not in (select rr.room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date);
	`

	row, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		var r models.Room
		err := row.Scan(&r.ID, &r.RoomName)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, r)
	}

	if err = row.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}

// GetRoomByID retrieves a room by its ID from the database.
func (m *postgresDbRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	query := `
	  select id, room_name, created_at, updated_at 
	  from rooms 
		where id = $1
	`
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return room, err
	}

	return room, nil
}

func (m *postgresDbRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User

	query := `
	  select id, first_name, last_name, email, password, access_level, created_at, updated_at
	  from users
		where id = $1
		`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.AccessLevel, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (m *postgresDbRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	  update users
	  set first_name = $1, last_name = $2, email = $3, access_level = $4, updated_at = $5
		from users
	`
	_, err := m.DB.ExecContext(ctx, stmt, u.FirstName, u.LastName, u.Email, u.AccessLevel, time.Now())
	if err != nil {
		return err
	}

	return nil
}

// Authenticate checks if the user exists and if the password is correct.
func (m *postgresDbRepo) Authenticate(email, password string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	id := 0
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return id, "", errors.New("incorrect password")
	} else if err != nil {
		return id, "", err
	}

	return id, hashedPassword, nil
}

func (m *postgresDbRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `
	  select r.id, r.first_name, r.last_name, r.email, r.phone_number, r.start_date, r.end_date,
	         r.room_id, r.created_at, r.updated_at,
	         rm.id, rm.room_name
	  from reservations r
		left join rooms rm on (r.room_id = rm.id)
		order by r.start_date asc
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.PhoneNumber,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

func (m *postgresDbRepo) NewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `
	  select r.id, r.first_name, r.last_name, r.email, r.phone_number, r.start_date, r.end_date,
	         r.room_id, r.created_at, r.updated_at, r.processed,
	         rm.id, rm.room_name
	  from reservations r
		left join rooms rm on (r.room_id = rm.id)
		where r.processed = false
		order by r.start_date asc
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.PhoneNumber,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

func (m *postgresDbRepo) GetReservationByID(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservation models.Reservation

	query := `
	  select r.id, r.first_name, r.last_name, r.email,
	         r.phone_number, r.start_date, r.end_date, r.room_id, 
	         r.created_at, r.updated_at, r.processed,
	         rm.id, rm.room_name
	  from reservations r
		left join rooms rm on (r.room_id = rm.id)
		where r.id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&reservation.ID,
		&reservation.FirstName,
		&reservation.LastName,
		&reservation.Email,
		&reservation.PhoneNumber,
		&reservation.StartDate,
		&reservation.EndDate,
		&reservation.RoomID,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
		&reservation.Processed,
		&reservation.Room.ID,
		&reservation.Room.RoomName,
	)
	if err != nil {
		return reservation, err
	}

	return reservation, nil
}

func (m *postgresDbRepo) UpdateReservation(res models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	  update reservations
	  set first_name = $1, last_name = $2, email = $3, phone_number = $4, updated_at = $5
		where id = $6
	`
	_, err := m.DB.ExecContext(ctx, stmt, res.FirstName, res.LastName, res.Email, res.PhoneNumber, time.Now(), res.ID)
	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDbRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	  delete from reservations
		where id = $1
	`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDbRepo) UpdateProcessedForReservation(id int, processed bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	  update reservations
	  set processed = $1, updated_at = $2
		where id = $3
	`

	_, err := m.DB.ExecContext(ctx, query, processed, time.Now(), id)
	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDbRepo) AllRooms() ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `
	  select id, room_name, created_at, updated_at
	  from rooms
		order by room_name
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return rooms, err
	}
	defer rows.Close()

	for rows.Next() {
		var rm models.Room
		err := rows.Scan(&rm.ID, &rm.RoomName, &rm.CreatedAt, &rm.UpdatedAt)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, rm)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

func (m *postgresDbRepo) GetRestrictionsForRoomByDates(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var restrictions []models.RoomRestriction

	query := `
	  select id, start_date, end_date, room_id, coalesce(reservation_id, 0), restriction_id, created_at, updated_at
	  from room_restrictions
		where $1 < end_date and $2 > start_date and room_id = $3
	`

	rows, err := m.DB.QueryContext(ctx, query, start, end, roomID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var r models.RoomRestriction
		err := rows.Scan(
			&r.ID,
			&r.StartDate,
			&r.EndDate,
			&r.RoomID,
			&r.ReservationID,
			&r.RestrictionID,
			&r.CreatedAt,
			&r.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		restrictions = append(restrictions, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return restrictions, nil
}

func (m *postgresDbRepo) InsertBlockForRoom(id int, start time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	  insert into room_restrictions (start_date, end_date, room_id, restriction_id, created_at, updated_at)
		values($1, $2, $3, $4, $5, $6)
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		start,
		start.AddDate(0, 0, 1), // Assuming a block lasts one day
		id,
		2,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

// Delete block by id
func (m *postgresDbRepo) DeleteBlockByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	  delete from room_restrictions
		where id = $1
	`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
