package handlers

import (
	"net/http"

	"github.com/nelsonmarro/bookings/config"
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
	remoteAddr := r.RemoteAddr
	h.app.Session.Put(r.Context(), "remote_ip", remoteAddr)

	home := templates.HomePage()
	err := home.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
