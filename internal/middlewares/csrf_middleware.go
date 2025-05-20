package middlewares

import (
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/nelsonmarro/bookings/config"
)

func CSRFMiddleware(next http.Handler, app *config.AppConfig) http.Handler {
	csrfHandler := nosurf.New(next) // nosurf wraps the *actual* next handler.

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/", // Ensure this path covers all parts of your app that need protection
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode, // Good default
	})
	return csrfHandler
}
