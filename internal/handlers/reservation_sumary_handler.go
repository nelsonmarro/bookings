package handlers

import (
	"net/http"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/models"
	"github.com/nelsonmarro/bookings/templates"
)

type ReservationSummaryHandler struct {
	app *config.AppConfig
}

func NewReservationSummaryHandler(app *config.AppConfig) *ReservationSummaryHandler {
	return &ReservationSummaryHandler{
		app: app,
	}
}

func (h *ReservationSummaryHandler) Get(w http.ResponseWriter, r *http.Request) {
	sessionData := h.app.Session.Get(r.Context(), "reservation")
	reservation, ok := sessionData.(models.Reservation)
	h.app.Session.Remove(r.Context(), "reservation")

	if !ok {
		h.app.Session.Put(r.Context(), "error", "No se encontraron datos para el resumen de la reservación.")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	summaryVm := templates.NewReservationSummaryPageVM()
	summaryVm.Reservation = reservation // Poblar con los datos de la sesión.
	summaryVm.StartDate = sd
	summaryVm.EndDate = ed
	summaryPage := templates.ReservationSumary(summaryVm)
	err := summaryPage.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
