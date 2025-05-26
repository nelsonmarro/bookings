package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/justinas/nosurf"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/helpers"
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
	vm := templates.NewConfirmReservationPageVM()
	vm.CSRFToken = nosurf.Token(r)
	reservation := templates.ConfirmReservationPage(vm)
	err := reservation.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h *ConfirmReservationHandler) Post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	sd := r.FormValue("start_date")
	ed := r.FormValue("end_date")

	dateLayout := "2006-01-02"

	startDate, err := time.Parse(dateLayout, sd)
	if err != nil {
		helpers.ServerError(w, err)
	}
	endDate, err := time.Parse(dateLayout, ed)
	if err != nil {
		helpers.ServerError(w, err)
	}

	roomID, err := strconv.Atoi(r.FormValue("room_id"))

	reservation := models.Reservation{
		FirstName:   r.FormValue("first_name"),
		LastName:    r.FormValue("last_name"),
		Email:       r.FormValue("email"),
		PhoneNumber: r.FormValue("phone_number"),
		StartDate:   startDate,
		EndDate:     endDate,
		RoomID:      roomID,
	}

	vm := templates.NewConfirmReservationPageVM()
	form := models.NewForm(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")

	if !form.Valid() {
		vm.Reservation = reservation
		vm.Form = form

		vm.CSRFToken = nosurf.Token(r)
		reservation := templates.ConfirmReservationPage(vm)
		err := reservation.Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
	}

	newRecordID, err := h.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomRestriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newRecordID,
		RestrictionID: 1,
	}

	err = h.DB.InsertRoomRestriction(roomRestriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	summaryVm := templates.NewReservationSummaryPageVM()
	summaryVm.Reservation = reservation
	h.app.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation/summary", http.StatusSeeOther)
}
