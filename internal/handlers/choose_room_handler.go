package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/helpers"
	"github.com/nelsonmarro/bookings/internal/models"
	"github.com/nelsonmarro/bookings/internal/repository"
)

type ChooseRoomHandler struct {
	app *config.AppConfig
	DB  repository.DataBaseRepo
}

func NewChooseRoomHandler(app *config.AppConfig, dbrepo repository.DataBaseRepo) *ChooseRoomHandler {
	return &ChooseRoomHandler{
		app: app,
		DB:  dbrepo,
	}
}

func (h *ChooseRoomHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	roomID, err := strconv.Atoi(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := h.app.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}

	res.RoomID = roomID

	h.app.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/reservation/confirmation", http.StatusSeeOther)
}
