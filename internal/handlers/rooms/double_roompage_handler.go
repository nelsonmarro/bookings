package rooms

import (
	"net/http"

	"github.com/justinas/nosurf"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/helpers"
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

func (h *DoubleRoomHandler) Get(w http.ResponseWriter, r *http.Request) {
	vm := rooms.NewDoubleRoomPageVM(nosurf.Token(r))

	doubleRoom := rooms.DoubleRoomPage(vm)
	err := doubleRoom.Render(r.Context(), w)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
}
