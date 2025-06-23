package admin

import (
	"net/http"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/repository"
	"github.com/nelsonmarro/bookings/templates/admin"
)

type AdminHandler struct {
	app *config.AppConfig
	DB  repository.DataBaseRepo
}

func NewAdminHandler(app *config.AppConfig, dbrepo repository.DataBaseRepo) *AdminHandler {
	return &AdminHandler{
		app: app,
		DB:  dbrepo,
	}
}

func (h *AdminHandler) GetAdminDashboard(w http.ResponseWriter, r *http.Request) {
	adminDashboard := admin.AdminDashboard()
	err := adminDashboard.Render(r.Context(), w)
	if err != nil {
		h.app.ErrorLog.Printf("error rendering admin dashboard: %v", err)
		return
	}
}
