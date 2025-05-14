package middlewares

import (
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/nelsonmarro/bookings/config"
)

// CSRFMiddleware is designed to be used with an anonymous function in mux.Use, like:
// mux.Use(func(next http.Handler) http.Handler {
//
//	    return middlewares.CSRFMiddleware(next, app)
//	})
//
// The 'next' parameter here is the actual next handler in the chain for a given route.
func CSRFMiddleware(next http.Handler, app *config.AppConfig) http.Handler {
	// Create and configure the nosurf handler. This happens when the middleware
	// chain is being built for a specific route or group of routes.
	csrfHandler := nosurf.New(next) // nosurf wraps the *actual* next handler.

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/", // Ensure this path covers all parts of your app that need protection
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode, // Good default
		// MaxAge:   3600 * 12, // Optional: Set a cookie lifetime (e.g., 12 hours)
		// Consider aligning MaxAge with your session cookie's MaxAge if applicable
	})

	// 1. On GET requests: Generate a token (if needed) and make it available via nosurf.Token(r).
	// 2. On POST requests: Check the submitted token against the one stored/derived from the cookie.
	// 3. If checks pass, call the `next` handler. If not, call the failure handler.
	return csrfHandler
}
