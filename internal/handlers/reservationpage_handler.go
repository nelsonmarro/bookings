package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/justinas/nosurf"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/models"
	"github.com/nelsonmarro/bookings/templates"
)

const htmlDateLayout = "2006-01-02"

type ReservationpageHandler struct {
	app *config.AppConfig
}

func NewReservationpageHandler(app *config.AppConfig) *ReservationpageHandler {
	return &ReservationpageHandler{
		app: app,
	}
}

func (h *ReservationpageHandler) Get(w http.ResponseWriter, r *http.Request) {
	vm := templates.NewReservationPageVM(nosurf.Token(r))

	reservation := templates.ReservationPage(vm)
	err := reservation.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h *ReservationpageHandler) Post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	vm := templates.NewReservationPageVM("")

	startDateStr := r.FormValue("startdate")
	endDateStr := r.FormValue("enddate")

	var parsedStartDate, parsedEndDate time.Time
	isValidStartDate := false
	isValidEndDate := false

	// Process Start Date
	if strings.TrimSpace(startDateStr) == "" {
		vm.Form.Errors.Add("startdate", "Start date is required")
	} else {
		t, parseErr := time.Parse(htmlDateLayout, startDateStr)
		if parseErr != nil {
			vm.Form.Errors.Add("startdate", "Invalid start date format. Please select a valid date.")
		} else {
			vm.StartDate = t
			parsedStartDate = t
			isValidStartDate = true
		}
	}

	// --- Process End Date ---
	if strings.TrimSpace(endDateStr) == "" {
		vm.Form.Errors.Add("enddate", "End date is required.")
	} else {
		t, parseErr := time.Parse(htmlDateLayout, endDateStr)
		if parseErr != nil {
			vm.Form.Errors.Add("enddate", "Invalid end date format. Please select a valid date.")
		} else {
			vm.EndDate = t // Set for re-populating the form
			parsedEndDate = t
			isValidEndDate = true
		}
	}

	// --- Cross-Field Validation ---
	if isValidStartDate && isValidEndDate {
		if parsedStartDate.After(parsedEndDate) {
			vm.Form.Errors.Add("enddate", "End date must be after start date.")
		}
	}

	// TODO:Fix validation logic
	if len(vm.Form.Errors) > 0 {
		// Render the form with errors
		w.WriteHeader(http.StatusBadRequest)
		vm.CSRFToken = nosurf.Token(r)
		reservation := templates.ReservationPage(vm)
		err := reservation.Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
		return
	}
	// --- All Validations Passed ---
	h.app.Session.Put(r.Context(), "flash_success", "Your availability check was successful!") // Example flash message
	http.Redirect(w, r, "/reservation/confirmation", http.StatusSeeOther)                      // Redirect back or to a summary page
}

func (h *ReservationpageHandler) PostJson(w http.ResponseWriter, r *http.Request) {
	resp := models.JsonResponse{
		Ok:      true,
		Message: "Reservation request received",
	}

	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		log.Println("Error marshalling JSON:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(out)
	if err != nil {
		log.Println("Error writing response:", err)
	}
}
