package rooms

import (
	"net/http"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/templates/rooms"
)

type DoubleRoomHandler struct {
	app *config.AppConfig
}

func NewDoubleRoomHandler(app *config.AppConfig) *DoubleRoomHandler {
	return &DoubleRoomHandler{
		app: app,
	}
}

func (h *DoubleRoomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	doubleRoom := rooms.DoubleRoomPage()
	err := doubleRoom.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
