package rooms

import (
	"net/http"

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
	vm := rooms.NewSingleRoomPageVM(nosurf.Token(r))

	singleRoom := rooms.SingleRoomPage(vm)
	err := singleRoom.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
