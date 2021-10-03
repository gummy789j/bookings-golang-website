package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/gummy789j/bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (this *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a restriction into the database
func (this *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {

	// Package context defines the Context type, which carries deadlines, cancellation signals, and other request-scoped values across API boundaries and between processes.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var newID int

	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	// You can execute a query, and return at most one row at the meantime by using QueryRowContext
	err := this.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
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

// InsertRoomRestriction inserts a room restriction into the database
func (this *postgresDBRepo) InsertRoomRestriction(rr models.RoomRestriction) error {

	// Package context defines the Context type, which carries deadlines, cancellation signals, and other request-scoped values across API boundaries and between processes.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at) 
			values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := this.DB.ExecContext(ctx, stmt,
		rr.StartDate,
		rr.EndDate,
		rr.RoomID,
		rr.ReservationID,
		rr.RestrictionID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityByDates return true if availability exists, and false if no availability exists
func (this *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var numRow int

	stmt := `select count(id) from room_restrictions rr where $1 < end_date and $2 > start_date and room_id = $3`

	err := this.DB.QueryRowContext(ctx, stmt, start, end, roomID).Scan(&numRow)
	if err != nil {
		return false, err
	}

	if numRow == 0 {
		return true, nil
	}

	return false, nil
}

// SearchAvailabilityForAllRooms return a slice of available rooms, if any; for given date range
func (this *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var rooms []models.Room

	query := `select r.id, r.room_name 
			from rooms r 
			where r.id not in (select rr.room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date)`

	rows, err := this.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {

		var room models.Room

		err := rows.Scan(
			&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return rooms, err
		}

		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

// GetRoomByID gets a room by id
func (this *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var room models.Room

	query := "select id, room_name, created_at, updated_at from rooms r where id = $1"

	err := this.DB.QueryRowContext(ctx, query, id).Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return room, err
	}

	return room, nil

}

// GetUserByID gets a user by id
func (this *postgresDBRepo) GetUserByID(id int) (models.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var user models.User

	query := `select first_name, last_name, email, password, access_level, created_at, updated_at from users where id = $1`
	row := this.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(user.FirstName, user.LastName, user.Email, user.Password, user.AccessLevel, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}

// UpdateUser updates a user in the database
func (this *postgresDBRepo) UpdateUser(user models.User) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `update users set 
	first_name = $1,
	last_name = $2,
	email = $3, 
	access_level = $4, 
	updated_at = $5 
	where id = $6
	`
	_, err := this.DB.ExecContext(ctx, query, user.FirstName, user.LastName, user.Email, user.AccessLevel, time.Now(), user.ID)
	if err != nil {
		return err
	}

	return nil
}

// Authenticate authenticate a user
func (this *postgresDBRepo) Authenticate(email, password string) (int, string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var id int
	var hashedPassword string

	row := this.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil

}

// AllReservations returns a slice of all reservations
func (this *postgresDBRepo) AllReservations() ([]models.Reservation, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var reservations []models.Reservation

	query := `select r.*, rm.id, rm.room_name 
	from reservations r
	left join rooms rm on (r.room_id = rm.id)
	order by r.start_date asc
	`

	rows, err := this.DB.QueryContext(ctx, query)
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
			&i.Phone,
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

// AllNewReservations returns a slice of all reservations
func (this *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var reservations []models.Reservation

	query := `select r.*, rm.id, rm.room_name 
	from reservations r
	left join rooms rm on (r.room_id = rm.id)
	where r.processed = 0
	order by r.start_date asc
	`

	rows, err := this.DB.QueryContext(ctx, query)
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
			&i.Phone,
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

// GetReservationByID returns one reservation by ID
func (this *postgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var res models.Reservation

	query := `
	select r.*, rm.id, rm.room_name 
	from reservations r 
	left join rooms rm on (r.room_id = rm.id) 
	where r.id = $1
	`

	row := this.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Phone,
		&res.StartDate,
		&res.EndDate,
		&res.RoomID,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.Processed,
		&res.Room.ID,
		&res.Room.RoomName,
	)
	if err != nil {
		return res, err
	}

	return res, nil
}

// UpdateReservation updates a reservation in the database
func (this *postgresDBRepo) UpdateReservation(res models.Reservation) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `update reservations set 
	first_name = $1,
	last_name = $2,
	email = $3, 
	phone = $4, 
	updated_at = $5 
	where id = $6
	`
	_, err := this.DB.ExecContext(ctx, query, res.FirstName, res.LastName, res.Email, res.Phone, time.Now(), res.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteReservarion deletes reservation by id
func (this *postgresDBRepo) DeleteReservation(id int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `delete from reservations where id = $1`

	_, err := this.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateProcessedForReservation updates processed for a reservation by id
func (this *postgresDBRepo) UpdateProcessedForReservation(id, processed int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `update reservations set processed = $1 where id = $2`

	_, err := this.DB.ExecContext(ctx, query, processed, id)
	if err != nil {
		return err
	}

	return nil
}

// AllRooms returns all rooms
func (this *postgresDBRepo) AllRooms() ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `select * from rooms order by room_name`

	rows, err := this.DB.QueryContext(ctx, query)
	if err != nil {
		return rooms, err
	}

	defer rows.Close()

	for rows.Next() {

		var room models.Room

		err = rows.Scan(
			&room.ID,
			&room.RoomName,
			&room.CreatedAt,
			&room.UpdatedAt,
		)
		if err != nil {
			return rooms, err
		}

		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

// GetRestrictionsForRoomByDate returns roomRestriction by dates range
func (this *postgresDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var restrictions []models.RoomRestriction

	query := `select id, coalesce(reservation_id, 0), restriction_id, room_id, start_date, end_date 
		from room_restrictions 
		where $1 < end_date and $2 >= start_date and room_id = $3
		`

	rows, err := this.DB.QueryContext(ctx, query, start, end, roomID)
	if err != nil {
		return restrictions, err
	}

	defer rows.Close()

	for rows.Next() {

		var r models.RoomRestriction

		err = rows.Scan(
			&r.ID,
			&r.ReservationID,
			&r.RestrictionID,
			&r.RoomID,
			&r.StartDate,
			&r.EndDate,
		)
		if err != nil {
			return restrictions, err
		}

		restrictions = append(restrictions, r)
	}

	if err = rows.Err(); err != nil {
		return restrictions, err
	}

	return restrictions, nil
}

// InsertBlockForRoom inserts a room restriction for block
func (this *postgresDBRepo) InsertBlockForRoom(roomID int, startDate time.Time) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, restriction_id, created_at, updated_at) 
			values ($1, $2, $3, $4, $5, $6)`

	_, err := this.DB.ExecContext(ctx, stmt,
		startDate,
		startDate.AddDate(0, 0, 1),
		roomID,
		2,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteBlockByID deletes a room restriction for block
func (this *postgresDBRepo) DeleteBlockByID(id int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	stmt := `delete from room_restrictions where id = $1 `

	_, err := this.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}
