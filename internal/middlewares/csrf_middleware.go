package middlewares

import (
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/nelsonmarro/bookings/config"
)

func CSRFMiddleware(next http.Handler, app *config.AppConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csrfHandler := nosurf.New(next)
		csrfHandler.SetBaseCookie(http.Cookie{
			HttpOnly: true,
			Path:     "/",
			Secure:   app.InProduction,
			SameSite: http.SameSiteLaxMode,
		})

		csrfHandler.ServeHTTP(w, r)
	})
}
