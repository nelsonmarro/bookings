package handlers

import (
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/nelsonmarro/bookings/config"
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
	vm := &templates.ConfirmReservationPageVM{
		FormErrors: make(map[string]string),
		CSRFToken:  nosurf.Token(r), // Get the CSRF token from the request
	}
	reservation := templates.ConfirmReservationPage(vm)
	err := reservation.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
