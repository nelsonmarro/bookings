package web

import (
	"encoding/gob"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/driver"
	"github.com/nelsonmarro/bookings/internal/handlers"
	"github.com/nelsonmarro/bookings/internal/handlers/rooms"
	"github.com/nelsonmarro/bookings/internal/middlewares"
	"github.com/nelsonmarro/bookings/internal/models"
	"github.com/nelsonmarro/bookings/internal/repository/dbrepo"
)

func Routes(app *config.AppConfig, db *driver.DB) http.Handler {
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	dbrepo := dbrepo.NewPostgresRepo(db.SQL, app)

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

	homepageHandler := handlers.NewHomepageHandler(app)
	mux.Get("/", homepageHandler.Get)

	aboutpageHandler := handlers.NewAboutpageHandler(app)
	mux.Get("/about", aboutpageHandler.Get)

	contactpageHandler := handlers.NewContactpageHandler(app)
	mux.Get("/contact", contactpageHandler.Get)

	reservationPageHandler := handlers.NewReservationpageHandler(app, dbrepo)
	mux.Get("/reservation", reservationPageHandler.Get)
	mux.Post("/reservation", reservationPageHandler.Post)
	mux.Post("/reservation-json", reservationPageHandler.PostJson)

	singleRoomHandler := rooms.NewSingleRoomHandler(app)
	mux.Get("/rooms/single", singleRoomHandler.Get)
	mux.Post("/rooms/single", singleRoomHandler.Post)

	confirmReservationHandler := handlers.NewConfirmReservationHandler(app, dbrepo)
	mux.Get("/reservation/confirmation", confirmReservationHandler.Get)
	mux.Post("/reservation/confirmation", confirmReservationHandler.Post)

	reservationSummaryHandler := handlers.NewReservationSummaryHandler(app)
	mux.Get("/reservation/summary", reservationSummaryHandler.Get)

	doubleRoomHandler := rooms.NewDoubleRoomHandler(app)
	mux.Get("/rooms/double", doubleRoomHandler.Get)

	return mux
}
