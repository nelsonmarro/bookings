package middlewares

import (
	"net/http"

	"github.com/nelsonmarro/bookings/config"
)

func SessionLoad(next http.Handler, app *config.AppConfig) http.Handler {
	return app.Session.LoadAndSave(next)
}
