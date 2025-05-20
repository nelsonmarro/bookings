package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/handlers"
	"github.com/nelsonmarro/bookings/internal/handlers/rooms"
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
	mux.Handle("GET /contact", handlers.NewContactpageHandler(app))

	reservationPageHandler := handlers.NewReservationpageHandler(app)

	mux.Get("/reservation", reservationPageHandler.Get)
	mux.Post("/reservation", reservationPageHandler.Post)
	mux.Post("/reservation-json", reservationPageHandler.PostJson)

	singleRoomHandler := rooms.NewSingleRoomHandler(app)
	mux.Get("/rooms/single", singleRoomHandler.Get)
	mux.Post("/rooms/single", singleRoomHandler.Post)

	confirmReservationHandler := handlers.NewConfirmReservationHandler(app)
	mux.Get("/reservation/confirmation", confirmReservationHandler.Get)
	mux.Post("/reservation/confirmation", confirmReservationHandler.Post)

	reservationSummaryHandler := handlers.NewReservationSummaryHandler(app)
	mux.Get("/reservation/summary", reservationSummaryHandler.Get)

	mux.Handle("GET /rooms/double", rooms.NewDoubleRoomHandler(app))

	return mux
}
