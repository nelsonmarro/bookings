package handlers

import (
	"encoding/gob"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/handlers/rooms"
	"github.com/nelsonmarro/bookings/internal/middlewares"
	"github.com/nelsonmarro/bookings/internal/models"
)

func getRoutes() http.Handler {
	app := config.GetConfigInstance()
	gob.Register(models.Reservation{})

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	// mux.Use(func(next http.Handler) http.Handler {
	// 	return middlewares.CSRFMiddleware(next, app)
	// })
	mux.Use(func(next http.Handler) http.Handler {
		return middlewares.SessionLoad(next, app)
	})
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static/", fs))

	homepageHandler := NewHomepageHandler(app)
	mux.Get("/", homepageHandler.Get)

	aboutpageHandler := NewAboutpageHandler(app)
	mux.Get("/about", aboutpageHandler.Get)

	contactpageHandler := NewContactpageHandler(app)
	mux.Get("/contact", contactpageHandler.Get)

	reservationPageHandler := NewReservationpageHandler(app)
	mux.Get("/reservation", reservationPageHandler.Get)
	mux.Post("/reservation", reservationPageHandler.Post)
	mux.Post("/reservation-json", reservationPageHandler.PostJson)

	singleRoomHandler := rooms.NewSingleRoomHandler(app)
	mux.Get("/rooms/single", singleRoomHandler.Get)
	mux.Post("/rooms/single", singleRoomHandler.Post)

	confirmReservationHandler := NewConfirmReservationHandler(app)
	mux.Get("/reservation/confirmation", confirmReservationHandler.Get)
	mux.Post("/reservation/confirmation", confirmReservationHandler.Post)

	reservationSummaryHandler := NewReservationSummaryHandler(app)
	mux.Get("/reservation/summary", reservationSummaryHandler.Get)

	doubleRoomHandler := rooms.NewDoubleRoomHandler(app)
	mux.Get("/rooms/double", doubleRoomHandler.Get)

	return mux
}
