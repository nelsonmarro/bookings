package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/helpers"
	"github.com/nelsonmarro/bookings/internal/models"
	"github.com/nelsonmarro/bookings/internal/repository"
)

type BookRoomHandler struct {
	app *config.AppConfig
	DB  repository.DataBaseRepo
}

func NewBookRoomHandler(app *config.AppConfig, dbrepo repository.DataBaseRepo) *BookRoomHandler {
	return &BookRoomHandler{
		app: app,
		DB:  dbrepo,
	}
}

func (h *BookRoomHandler) Get(w http.ResponseWriter, r *http.Request) {
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sdate := r.URL.Query().Get("s")
	endate := r.URL.Query().Get("e")

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

	room, err := h.DB.GetRoomByID(roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var res models.Reservation
	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate
	res.Room.RoomName = room.RoomName

	h.app.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/reservation/confirmation", http.StatusSeeOther)
}
