package admin

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/justinas/nosurf"
	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/helpers"
	"github.com/nelsonmarro/bookings/internal/models"
	"github.com/nelsonmarro/bookings/internal/repository"
	"github.com/nelsonmarro/bookings/templates/admin"
)

type AdminReservationsHandler struct {
	app *config.AppConfig
	DB  repository.DataBaseRepo
}

func NewAdminReservationsHandler(app *config.AppConfig, db repository.DataBaseRepo) *AdminReservationsHandler {
	return &AdminReservationsHandler{
		app: app,
		DB:  db,
	}
}

func (h *AdminReservationsHandler) GetNewReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := h.DB.NewReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	vm := admin.NewAdminNewReservationsVM(reservations)
	messageType, message := models.GetSessionMessage(r.Context())
	vm.MessageType = messageType
	vm.Message = message

	newReservationsPage := admin.AdminNewReservationsPage(vm)
	err = newReservationsPage.Render(r.Context(), w)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
}

func (h *AdminReservationsHandler) GetAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := h.DB.AllReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	vm := admin.NewAdminllReservationsVM(reservations)
	messageType, message := models.GetSessionMessage(r.Context())
	vm.MessageType = messageType
	vm.Message = message

	allReservationsPage := admin.AdminAllReservationsPage(vm)
	err = allReservationsPage.Render(r.Context(), w)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
}

func (h *AdminReservationsHandler) GetReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	// assume that there is no month/year specified in the URL
	now := time.Now()

	if r.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		month, _ := strconv.Atoi(r.URL.Query().Get("m"))
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")

	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")

	thisMonth := now.Format("01")
	thisMonthYear := now.Format("2006")

	// get the first and last day of the month
	currentYear, currentMonth, _ := now.Date()
	currentLoaction := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLoaction)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	rooms, err := h.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	vm := admin.NewAdminCalendarVM(
		nextMonth,
		nextMonthYear,
		lastMonth,
		lastMonthYear,
		thisMonth,
		thisMonthYear,
		now,
		lastOfMonth.Day(),
		rooms,
	)
	messageType, message := models.GetSessionMessage(r.Context())
	vm.MessageType = messageType
	vm.Message = message
	vm.CSRFToken = nosurf.Token(r)

	for _, room := range rooms {
		reservationsMap := make(map[string]int, len(rooms))
		blockMap := make(map[string]int, len(rooms))

		for d := firstOfMonth; !d.After(lastOfMonth); d = d.AddDate(0, 0, 1) {
			reservationsMap[d.Format("2006-01-02")] = 0
			blockMap[d.Format("2006-01-02")] = 0
		}

		// get all restrictions for the room
		restrictions, err := h.DB.GetRestrictionsForRoomByDates(room.ID, firstOfMonth, lastOfMonth)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		for _, restriction := range restrictions {
			if restriction.ReservationID > 0 {
				// its a reservation
				for d := restriction.StartDate; !d.After(restriction.EndDate); d = d.AddDate(0, 0, 1) {
					reservationsMap[d.Format("2006-01-02")] = restriction.ReservationID
				}
			} else {
				// its a block
				blockMap[restriction.StartDate.Format("2006-01-02")] = restriction.ID
			}
		}
		vm.ReservationMaps[fmt.Sprintf("reservation_map_%d", room.ID)] = reservationsMap
		vm.BlockMaps[fmt.Sprintf("block_map_%d", room.ID)] = blockMap

		h.app.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", room.ID), blockMap)
	}

	calendarPage := admin.AdminCalendarPage(vm)
	err = calendarPage.Render(r.Context(), w)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
}

func (h *AdminReservationsHandler) GetReservation(w http.ResponseWriter, r *http.Request) {
	// get url params
	src := chi.URLParam(r, "src")
	idStr := chi.URLParam(r, "id")

	m := r.URL.Query().Get("m")
	y := r.URL.Query().Get("y")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation, err := h.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	vm := admin.NewReservationDetailVM(reservation, src)
	vm.CSRFToken = nosurf.Token(r)

	vm.ResMonth = m
	vm.ResYear = y

	calendarPage := admin.AdminReservationDetail(vm)
	err = calendarPage.Render(r.Context(), w)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
}

func (h *AdminReservationsHandler) GetProcessReservation(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	src := chi.URLParam(r, "src")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	err = h.DB.UpdateProcessedForReservation(id, true)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m := r.URL.Query().Get("m")
	y := r.URL.Query().Get("y")
	redirectURL := fmt.Sprintf("/admin/reservations-%s", src)

	if m != "" && y != "" {
		redirectURL += fmt.Sprintf("?m=%s&y=%s", m, y)
	}

	h.app.Session.Put(r.Context(), "info", "Reservation marked as processed successfully!")
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (h *AdminReservationsHandler) GetDeleteReservation(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	src := chi.URLParam(r, "src")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	err = h.DB.DeleteReservation(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m := r.URL.Query().Get("m")
	y := r.URL.Query().Get("y")
	redirectURL := fmt.Sprintf("/admin/reservations-%s", src)

	if m != "" && y != "" {
		redirectURL += fmt.Sprintf("?m=%s&y=%s", m, y)
	}

	h.app.Session.Put(r.Context(), "info", "Reservation deleted successfully!")
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (h *AdminReservationsHandler) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := chi.URLParam(r, "src")
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, err := h.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.PhoneNumber = r.Form.Get("phone_number")

	err = h.DB.UpdateReservation(res)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m := r.URL.Query().Get("m")
	y := r.URL.Query().Get("y")
	redirectURL := fmt.Sprintf("/admin/reservations-%s", src)

	if m != "" && y != "" {
		redirectURL += fmt.Sprintf("?m=%s&y=%s", m, y)
	}

	h.app.Session.Put(r.Context(), "info", "Reservation updated successfully!")
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (h *AdminReservationsHandler) PostReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	year, _ := strconv.Atoi(r.Form.Get("y"))
	month, _ := strconv.Atoi(r.Form.Get("m"))

	// process blocks
	rooms, err := h.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := models.NewForm(r.PostForm)
	for _, room := range rooms {
		// Get the block map from the session.
		blockMap := h.app.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", room.ID)).(map[string]int)
		for name, value := range blockMap {
			if val, ok := blockMap[name]; ok {
				if val > 0 {
					if !form.Has(fmt.Sprintf("remove_block_%d_%s", room.ID, name), r) {
						// delete the restriction by id
						err = h.DB.DeleteBlockByID(value)
						if err != nil {
							helpers.ServerError(w, err)
							return
						}
					}
				}
			}
		}
	}

	// handle new blocks
	for name := range r.PostForm {
		if strings.HasPrefix(name, "add_block") {
			exploded := strings.Split(name, "_")
			roomID, _ := strconv.Atoi(exploded[2])

			t, _ := time.Parse("2006-01-02", exploded[3])
			// insert a new block
			err := h.DB.InsertBlockForRoom(roomID, t)
			if err != nil {
				helpers.ServerError(w, err)
				return
			}
		}
	}

	h.app.Session.Put(r.Context(), "info", "Changes Saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%d&m=%02d", year, month), http.StatusSeeOther)
}
