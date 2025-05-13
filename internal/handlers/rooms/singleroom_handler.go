package rooms

import (
	"net/http"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/templates/rooms"
)

type SingleRoomHandler struct {
	app *config.AppConfig
}

func NewSingleRoomHandler(app *config.AppConfig) *SingleRoomHandler {
	return &SingleRoomHandler{
		app: app,
	}
}

func (h *SingleRoomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	singleRoom := rooms.SingleRoomPage()
	err := singleRoom.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
