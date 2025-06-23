package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/justinas/nosurf"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/helpers"
	"github.com/nelsonmarro/bookings/internal/models"
	"github.com/nelsonmarro/bookings/internal/repository"
	"github.com/nelsonmarro/bookings/templates"
)

const htmlDateLayout = "2006-01-02"

type ReservationpageHandler struct {
	app *config.AppConfig
	DB  repository.DataBaseRepo
}

func NewReservationpageHandler(app *config.AppConfig, dbrepo repository.DataBaseRepo) *ReservationpageHandler {
	return &ReservationpageHandler{
		app: app,
		DB:  dbrepo,
	}
}

func (h *ReservationpageHandler) Get(w http.ResponseWriter, r *http.Request) {
	vm := templates.NewReservationPageVM(nosurf.Token(r))
	messageType, message := models.GetSessionMessage(r.Context())
	vm.MessageType = messageType
	vm.Message = message

	reservation := templates.ReservationPage(vm)
	err := reservation.Render(r.Context(), w)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
}

func (h *ReservationpageHandler) Post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	vm := templates.NewReservationPageVM("")

	startDateStr := r.FormValue("start_date")
	endDateStr := r.FormValue("end_date")

	var parsedStartDate, parsedEndDate time.Time

	// Process Start Date
	if strings.TrimSpace(startDateStr) == "" {
		vm.Form.Errors.Add("start_date", "Start date is required")
	} else {
		t, parseErr := time.Parse(htmlDateLayout, startDateStr)
		if parseErr != nil {
			vm.Form.Errors.Add("start_date", "Invalid start date format. Please select a valid date.")
		} else {
			vm.StartDate = t
			parsedStartDate = t
		}
	}

	// --- Process End Date ---
	if strings.TrimSpace(endDateStr) == "" {
		vm.Form.Errors.Add("end_date", "End date is required.")
	} else {
		t, parseErr := time.Parse(htmlDateLayout, endDateStr)
		if parseErr != nil {
			vm.Form.Errors.Add("end_date", "Invalid end date format. Please select a valid date.")
		} else {
			vm.EndDate = t // Set for re-populating the form
			parsedEndDate = t
		}
	}

	if len(vm.Form.Errors) > 0 {
		// Render the form with errors
		w.WriteHeader(http.StatusBadRequest)
		vm.CSRFToken = nosurf.Token(r)
		reservation := templates.ReservationPage(vm)
		err := reservation.Render(r.Context(), w)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		return
	}

	rooms, err := h.DB.SearchAvailabilityForAllRooms(parsedStartDate, parsedEndDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	for _, room := range rooms {
		h.app.InfoLog.Printf("Available room: %s (ID: %d)", room.RoomName, room.ID)
	}

	if len(rooms) == 0 {
		// no available rooms
		h.app.Session.Put(r.Context(), "error", "No rooms available for the selected dates.")
		http.Redirect(w, r, "/reservation", http.StatusSeeOther)
		return
	}

	res := models.Reservation{
		StartDate: vm.StartDate,
		EndDate:   vm.EndDate,
	}
	h.app.Session.Put(r.Context(), "reservation", res)

	chooseRoomVm := templates.NewChooseRoomPageVM(rooms)
	chooseRoomPage := templates.ChooseRoomPage(chooseRoomVm)
	err = chooseRoomPage.Render(r.Context(), w)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
}

func (h *ReservationpageHandler) PostJson(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	sdate := r.FormValue("start_date")
	endate := r.FormValue("end_date")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sdate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, endate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomID, _ := strconv.Atoi(r.FormValue("room_id"))

	available, err := h.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	resp := models.JsonResponse{
		Ok:        available,
		Message:   "",
		StartDate: sdate,
		EndDate:   endate,
		RoomID:    strconv.Itoa(roomID),
	}

	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(out)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
}
