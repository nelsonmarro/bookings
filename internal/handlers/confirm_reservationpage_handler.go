package handlers

import (
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/models"
	"github.com/nelsonmarro/bookings/templates"
)

type ConfirmReservationHandler struct {
	app *config.AppConfig
}

func NewConfirmReservationHandler(app *config.AppConfig) *ConfirmReservationHandler {
	return &ConfirmReservationHandler{
		app: app,
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

	reservation := models.Reservation{
		FirstName:   r.FormValue("first_name"),
		LastName:    r.FormValue("last_name"),
		Email:       r.FormValue("email"),
		PhoneNumber: r.FormValue("phone_number"),
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

	summaryVm := templates.NewReservationSummaryPageVM()
	summaryVm.Reservation = reservation
	h.app.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation/summary", http.StatusSeeOther)
}
