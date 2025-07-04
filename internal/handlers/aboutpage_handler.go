package handlers

import (
	"net/http"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/templates"
)

type AboutpageHandler struct {
	app *config.AppConfig
}

func NewAboutpageHandler(app *config.AppConfig) *AboutpageHandler {
	return &AboutpageHandler{
		app: app,
	}
}

func (h *AboutpageHandler) Get(w http.ResponseWriter, r *http.Request) {
	about := templates.AboutePage()
	err := about.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
