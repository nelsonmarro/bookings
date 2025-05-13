package handlers

import (
	"net/http"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/templates"
)

type ContactpageHandler struct {
	app *config.AppConfig
}

func NewContactpageHandler(app *config.AppConfig) *ContactpageHandler {
	return &ContactpageHandler{
		app: app,
	}
}

func (h *ContactpageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contact := templates.ContactPage()
	err := contact.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
