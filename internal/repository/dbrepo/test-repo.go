package dbrepo

import (
	"errors"
	"time"

	"github.com/gummy789j/bookings/internal/models"
)

func (this *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a restriction into the database
func (this *testDBRepo) InsertReservation(res models.Reservation) (int, error) {

	if res.RoomID == 1000 {
		return 0, errors.New("Some error!")
	}
	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (this *testDBRepo) InsertRoomRestriction(rr models.RoomRestriction) error {

	if rr.RoomID == 1001 {
		return errors.New("Some error!")
	}

	return nil
}

// SearchAvailabilityByDates return true if availability exists, and false if no availability exists
func (this *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {

	if roomID == 1000 {
		return false, errors.New("Some error!")
	}
	return false, nil
}

// SearchAvailabilityForAllRooms return a slice of available rooms, if any; for given date range
func (this *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil
}

// GetRoomByID gets a room by id
func (this *testDBRepo) GetRoomByID(id int) (models.Room, error) {

	var room models.Room

	if id == 1 {
		return room, errors.New("Some error!")
	}

	return room, nil

}

// GetUserByID gets a user by id
func (this *testDBRepo) GetUserByID(id int) (models.User, error) {

	var user models.User

	return user, nil
}

// UpdateUser updates a user in the database
func (this *testDBRepo) UpdateUser(user models.User) error {

	return nil

}

// Authenticate authenticate a user
func (this *testDBRepo) Authenticate(email, password string) (int, string, error) {
	return 0, "", nil
}

// AllReservations returns a slice of all reservations
func (this *testDBRepo) AllReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation

	return reservations, nil
}

// AllNewReservations returns a slice of all reservations
func (this *testDBRepo) AllNewReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation

	return reservations, nil
}

// GetReservationByID returns one reservation by ID
func (this *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {

	var res models.Reservation

	return res, nil
}

// UpdateReservation updates a reservation in the database
func (this *testDBRepo) UpdateReservation(res models.Reservation) error {

	return nil
}

// DeleteReservarion deletes reservation by id
func (this *testDBRepo) DeleteReservation(id int) error {

	return nil
}

// UpdateProcessedForReservation updates processed for a reservation by id
func (this *testDBRepo) UpdateProcessedForReservation(id, processed int) error {

	return nil
}

// AllRooms returns all rooms
func (this *testDBRepo) AllRooms() ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil
}

// GetRestrictionsForRoomByDate returns roomRestriction by dates range
func (this *testDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {

	var restrictions []models.RoomRestriction

	return restrictions, nil
}

// InsertBlockForRoom inserts a room restriction for block
func (this *testDBRepo) InsertBlockForRoom(roomID int, startDate time.Time) error {

	return nil
}

// DeleteBlockByID deletes a room restriction for block
func (this *testDBRepo) DeleteBlockByID(id int) error {

	return nil
}
