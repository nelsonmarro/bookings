package rooms

import (
	"net/http"
	"strings"
	"time"

	"github.com/justinas/nosurf"
	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/templates/rooms"
)

const htmlDateLayout = "2006-01-02"

type SingleRoomHandler struct {
	app *config.AppConfig
}

func NewSingleRoomHandler(app *config.AppConfig) *SingleRoomHandler {
	return &SingleRoomHandler{
		app: app,
	}
}

func (h *SingleRoomHandler) Get(w http.ResponseWriter, r *http.Request) {
	vm := &rooms.SingleRoomPageVM{
		FormErrors: make(map[string]string),
		CSRFToken:  nosurf.Token(r), // Get the CSRF token from the request
	}

	singleRoom := rooms.SingleRoomPage(vm)
	err := singleRoom.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h *SingleRoomHandler) Post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	vm := &rooms.SingleRoomPageVM{
		FormErrors: make(map[string]string),
	}

	startDateStr := r.FormValue("startdate")
	endDateStr := r.FormValue("enddate")

	var parsedStartDate, parsedEndDate time.Time
	isValidStartDate := false
	isValidEndDate := false

	// Process Start Date
	if strings.TrimSpace(startDateStr) == "" {
		vm.FormErrors["startdate"] = "Start date is required"
	} else {
		t, parseErr := time.Parse(htmlDateLayout, startDateStr)
		if parseErr != nil {
			vm.FormErrors["startdate"] = "Invalid start date format. Please select a valid date."
		} else {
			vm.StartDate = t
			parsedStartDate = t
			isValidStartDate = true
		}
	}

	// --- Process End Date ---
	if strings.TrimSpace(endDateStr) == "" {
		vm.FormErrors["enddate"] = "End date is required."
	} else {
		t, parseErr := time.Parse(htmlDateLayout, endDateStr)
		if parseErr != nil {
			vm.FormErrors["enddate"] = "Invalid end date format. Please select a valid date."
		} else {
			vm.EndDate = t // Set for re-populating the form
			parsedEndDate = t
			isValidEndDate = true
		}
	}

	// --- Cross-Field Validation ---
	if isValidStartDate && isValidEndDate {
		if parsedStartDate.After(parsedEndDate) {
			vm.FormErrors["enddate"] = "End date must be after start date."
		}
	}

	if len(vm.FormErrors) > 0 {
		// Render the form with errors
		w.WriteHeader(http.StatusBadRequest)
		vm.CSRFToken = nosurf.Token(r)
		singleRoom := rooms.SingleRoomPage(vm)
		err := singleRoom.Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
		return
	}
	// --- All Validations Passed ---
	// Here you would:
	// 1. Check availability using a service/repository with parsedStartDate, parsedEndDate.
	// 2. If not available, add to vm.FormErrors["general"] = "Not available" and re-render.
	// 3. If available, potentially create a pending reservation or proceed to payment.

	// For now, let's assume success and redirect (PRG pattern)
	h.app.Session.Put(r.Context(), "flash_success", "Your availability check was successful!") // Example flash message
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)                              // Redirect back or to a summary page
}
