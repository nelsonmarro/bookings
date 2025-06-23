package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/justinas/nosurf"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/email"
	"github.com/nelsonmarro/bookings/internal/models"
	"github.com/nelsonmarro/bookings/internal/repository"
	"github.com/nelsonmarro/bookings/templates"
)

type ConfirmReservationHandler struct {
	app *config.AppConfig
	DB  repository.DataBaseRepo
}

func NewConfirmReservationHandler(app *config.AppConfig, dbrepo repository.DataBaseRepo) *ConfirmReservationHandler {
	return &ConfirmReservationHandler{
		app: app,
		DB:  dbrepo,
	}
}

func (h *ConfirmReservationHandler) Get(w http.ResponseWriter, r *http.Request) {
	res, ok := h.app.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		h.app.Session.Put(r.Context(), "error", "Could not get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := h.DB.GetRoomByID(res.RoomID)
	if err != nil {
		h.app.Session.Put(r.Context(), "error", "Cannot find room for reservation")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res.Room = room

	h.app.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	vm := templates.NewConfirmReservationPageVM(res)

	vm.StartDate = sd
	vm.EndDate = ed
	vm.CSRFToken = nosurf.Token(r)
	reservation := templates.ConfirmReservationPage(vm)
	err = reservation.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h *ConfirmReservationHandler) Post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.app.Session.Put(r.Context(), "error", "Cannot parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := r.FormValue("start_date")
	ed := r.FormValue("end_date")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		h.app.Session.Put(r.Context(), "error", "Invalid start date format")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		h.app.Session.Put(r.Context(), "error", "Invalid end date format")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	roomID, err := strconv.Atoi(r.FormValue("room_id"))
	if err != nil {
		h.app.Session.Put(r.Context(), "error", "Invalid room ID")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res := models.Reservation{
		FirstName:   r.FormValue("first_name"),
		LastName:    r.FormValue("last_name"),
		PhoneNumber: r.FormValue("phone_number"),
		Email:       r.FormValue("email"),
		StartDate:   startDate,
		EndDate:     endDate,
		RoomID:      roomID,
	}

	form := models.NewForm(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")

	vm := templates.NewConfirmReservationPageVM(res)
	if !form.Valid() {
		vm.Form = form
		vm.CSRFToken = nosurf.Token(r)
		reservation := templates.ConfirmReservationPage(vm)
		err := reservation.Render(r.Context(), w)
		if err != nil {
			h.app.Session.Put(r.Context(), "error", "Cannot reneder page")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		return
	}

	newRecordID, err := h.DB.InsertReservation(res)
	if err != nil {
		h.app.Session.Put(r.Context(), "error", "Cannot insert reservation into database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	h.app.Session.Put(r.Context(), "reservation", res)

	roomRestriction := models.RoomRestriction{
		StartDate:     res.StartDate,
		EndDate:       res.EndDate,
		RoomID:        res.RoomID,
		ReservationID: newRecordID,
		RestrictionID: 1,
	}

	err = h.DB.InsertRoomRestriction(roomRestriction)
	if err != nil {
		h.app.Session.Put(r.Context(), "error", "Cannot insert room restriction into database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	htmlMessage := fmt.Sprintf(`
		  <string>Reservation Confirmation</string>
		  <br>
		  Dear %s: <br>
      This is confirmation of your reservation from %s to %s.<br>
		`, res.FirstName, res.StartDate.Format("2006-01-02"), res.EndDate.Format("2006-01-02"))

	msg := email.MailData{
		To:       res.Email,
		From:     "me@here.com",
		Content:  htmlMessage,
		Subject:  "Reservation Confirmation",
		Template: "basic.html",
	}

	h.app.MailChan <- msg

	// send message to property owner
	htmlMessage = fmt.Sprintf(`
		  <string>Reservation Notification</string>
		  <br>
		  a new reservation has been made from %s to %s for %s.<br>
		`, res.StartDate.Format("2006-01-02"), res.EndDate.Format("2006-01-02"), res.Room.RoomName)

	msg = email.MailData{
		To:       "me@here.com",
		From:     "me@here.com",
		Subject:  "Reservation Notification",
		Content:  htmlMessage,
		Template: "basic.html",
	}

	h.app.MailChan <- msg

	http.Redirect(w, r, "/reservation/summary", http.StatusSeeOther)
}
