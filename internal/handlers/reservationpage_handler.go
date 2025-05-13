package handlers

import (
	"net/http"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/templates"
)

type ReservationpageHandler struct {
	app *config.AppConfig
}

func NewReservationpageHandler(app *config.AppConfig) *ReservationpageHandler {
	return &ReservationpageHandler{
		app: app,
	}
}

func (h *ReservationpageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reservation := templates.ReservationPage()
	err := reservation.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
