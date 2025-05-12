package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/handlers"
	"github.com/nelsonmarro/bookings/internal/middlewares"
)

func Routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(func(next http.Handler) http.Handler {
		return middlewares.CSRFMiddleware(next, app)
	})
	mux.Use(func(next http.Handler) http.Handler {
		return middlewares.SessionLoad(next, app)
	})

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static/", fs))

	mux.Handle("GET /", handlers.NewHomepageHandler(app))
	mux.Handle("GET /about", handlers.NewAboutpageHandler(app))

	return mux
}
