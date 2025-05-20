package handlers

import (
	"net/http"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/models"
	"github.com/nelsonmarro/bookings/templates"
)

type HomepageHandler struct {
	app *config.AppConfig
}

func NewHomepageHandler(app *config.AppConfig) *HomepageHandler {
	return &HomepageHandler{
		app: app,
	}
}

func (h *HomepageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vm := templates.NewHomePageVM()

	messageType, message := models.GetSessionMessage(h.app, r)
	vm.MessageType = messageType
	vm.Message = message

	home := templates.HomePage(vm)
	err := home.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
