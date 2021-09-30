package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/gummy789j/bookings/internal/config"
	"github.com/gummy789j/bookings/internal/driver"
	"github.com/gummy789j/bookings/internal/forms"
	"github.com/gummy789j/bookings/internal/helpers"
	"github.com/gummy789j/bookings/internal/models"
	"github.com/gummy789j/bookings/internal/render"
	"github.com/gummy789j/bookings/internal/repository"
	"github.com/gummy789j/bookings/internal/repository/dbrepo"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewTestRepo creates a new repository
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestRepo(a),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

// Home
func (this *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr

	this.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About
func (this *Repository) About(w http.ResponseWriter, r *http.Request) {

	stringMap := map[string]string{
		"test": "Hello, again",
	}

	remoteIP := this.App.Session.GetString(r.Context(), "remote_ip")

	stringMap["remote_ip"] = remoteIP

	render.Template(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// search-availability
func (this *Repository) Availability(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// Post search-availability
func (this *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := this.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	for _, val := range rooms {
		this.App.InfoLog.Println("Room", val.ID, val.RoomName)
	}

	if len(rooms) == 0 {
		this.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	data["rooms"] = rooms

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	this.App.Session.Put(r.Context(), "reservation", res)
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// Json search-availability
func (this *Repository) JsonAvailability(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error",
		}

		// indent 縮進
		out, _ := json.MarshalIndent(resp, "", "     ")

		//log.Println(string(out))
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)

	endDate, _ := time.Parse(layout, ed)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, err := this.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Error connecting to database",
		}

		// indent 縮進
		out, _ := json.MarshalIndent(resp, "", "     ")

		//log.Println(string(out))
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	resp := jsonResponse{
		OK:        available,
		Message:   "",
		RoomID:    strconv.Itoa(roomID),
		StartDate: sd,
		EndDate:   ed,
	}

	// indent 縮進
	out, _ := json.MarshalIndent(resp, "", "     ")

	//log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Make-reservation
func (this *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	res, ok := this.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		// helpers.ServerError(w, errors.New("Cannot get reservation"))
		// return
		this.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := this.DB.GetRoomByID(res.RoomID)
	if err != nil {
		// helpers.ServerError(w, err)
		// return
		this.App.Session.Put(r.Context(), "error", "Can't find room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room = room

	this.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})

	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})

}

// Post Make-reservation
func (this *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {

	// res, ok := this.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	// if !ok {
	// 	this.App.Session.Put(r.Context(), "error", "Can't parse form!")
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }

	err := r.ParseForm()
	if err != nil {
		this.App.Session.Put(r.Context(), "error", "Can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		this.App.Session.Put(r.Context(), "error", "Can't parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		this.App.Session.Put(r.Context(), "error", "Can't parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		this.App.Session.Put(r.Context(), "error", "Invalid data")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")

	form.MinLength("first_name", 3)

	form.IsEmail("email")

	if !form.Valid() {

		data := make(map[string]interface{})
		data["reservation"] = reservation
		http.Error(w, "my own error message", http.StatusSeeOther)
		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newReservationID, err := this.DB.InsertReservation(reservation)
	if err != nil {
		this.App.Session.Put(r.Context(), "error", "Can't insert reservation into database!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = this.DB.InsertRoomRestriction(restriction)
	if err != nil {
		this.App.Session.Put(r.Context(), "error", "Can't insert room restriction!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// send notification - first to guest)
	htmlMessage := fmt.Sprintf(
		`
		<strong>Reservation Confirmation</strong> <br>
		Dear %s: <br>
		This is confirm your reservation from %s to %s.
		`,
		reservation.FirstName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	msg := models.MailData{
		To:       reservation.Email,
		From:     "linshotel@hotel.com",
		Subject:  "Reservation Confirmation",
		Content:  htmlMessage,
		Template: "basic.html",
	}

	this.App.MailChan <- msg

	this.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Generals
func (this *Repository) Generals(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors
func (this *Repository) Majors(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Contact
func (this *Repository) Contact(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// ReservationSummary
func (this *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {

	// it need to assert to models.Reservation type
	reservation, ok := this.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		this.App.ErrorLog.Println("Can't get error from session")
		this.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	this.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	stringMap := make(map[string]string)

	stringMap["start_date"] = reservation.StartDate.Format("2006-01-02")
	stringMap["end_date"] = reservation.EndDate.Format("2006-01-02")

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// ChooseRoom
func (this *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := this.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}

	res.RoomID = roomID

	this.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// BookRoom takes URL parameters, builds a seesional variable, and takes user to make reservation screen
func (this *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {

	ID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	var res models.Reservation

	res.RoomID = ID
	room, err := this.DB.GetRoomByID(ID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.Room = room

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.StartDate = startDate
	res.EndDate = endDate

	this.App.Session.Put(r.Context(), "reservation", res)

	//log.Println(ID, startDate, endDate)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

// ShowLogin
func (this *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

//PostShowLogin handles logging the user in
func (this *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {

	// 使user每次登入登出時，都會使用新的session id
	// 避免 fixation attack
	_ = this.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	form := forms.New(r.PostForm)

	form.IsEmail("email")
	form.Required("email", "password")
	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := this.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		this.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	this.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	this.App.Session.Put(r.Context(), "user_id", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout logs a user out
func (this *Repository) Logout(w http.ResponseWriter, r *http.Request) {

	_ = this.App.Session.Destroy(r.Context())
	_ = this.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (this *Repository) AdminDashBoard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}

func (this *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "admin-new-reservations.page.tmpl", &models.TemplateData{})

}

func (this *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "admin-all-reservations.page.tmpl", &models.TemplateData{})

}

func (this *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "admin-reservations-calendar.page.tmpl", &models.TemplateData{})

}
