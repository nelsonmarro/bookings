package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/driver"
	"github.com/nelsonmarro/bookings/internal/handlers"
	"github.com/nelsonmarro/bookings/internal/handlers/admin"
	"github.com/nelsonmarro/bookings/internal/handlers/rooms"
	"github.com/nelsonmarro/bookings/internal/middlewares"
	"github.com/nelsonmarro/bookings/internal/repository/dbrepo"
)

func Routes(app *config.AppConfig, db *driver.DB) http.Handler {
	dbrepo := dbrepo.NewPostgresRepo(db.SQL, app)

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(func(next http.Handler) http.Handler {
		return middlewares.CSRFMiddleware(next, app)
	})
	mux.Use(app.Session.LoadAndSave)

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

	chooseRoomHandler := handlers.NewChooseRoomHandler(app, dbrepo)
	mux.Get("/choose-room/{id}", chooseRoomHandler.Get)

	bookRoomHandler := handlers.NewBookRoomHandler(app, dbrepo)
	mux.Get("/book-room", bookRoomHandler.Get)

	singleRoomHandler := rooms.NewSingleRoomHandler(app)
	mux.Get("/rooms/single", singleRoomHandler.Get)

	confirmReservationHandler := handlers.NewConfirmReservationHandler(app, dbrepo)
	mux.Get("/reservation/confirmation", confirmReservationHandler.Get)
	mux.Post("/reservation/confirmation", confirmReservationHandler.Post)

	reservationSummaryHandler := handlers.NewReservationSummaryHandler(app)
	mux.Get("/reservation/summary", reservationSummaryHandler.Get)

	doubleRoomHandler := rooms.NewDoubleRoomHandler(app)
	mux.Get("/rooms/double", doubleRoomHandler.Get)

	userHandler := handlers.NewUserHandler(app, dbrepo)
	mux.Get("/user/login", userHandler.GetLogin)
	mux.Post("/user/login", userHandler.PostLogin)
	mux.Get("/user/logout", userHandler.GetLogout)

	mux.Route("/admin", func(r chi.Router) {
		// r.Use(middlewares.Auth)

		adminHandler := admin.NewAdminHandler(app, dbrepo)
		r.Get("/dashboard", adminHandler.GetAdminDashboard)

		reservationHandler := admin.NewAdminReservationsHandler(app, dbrepo)
		r.Get("/reservations-new", reservationHandler.GetNewReservations)
		r.Get("/reservations-all", reservationHandler.GetAllReservations)
		r.Get("/reservations-calendar", reservationHandler.GetReservationsCalendar)
		r.Post("/reservations-calendar", reservationHandler.PostReservationsCalendar)
		r.Get("/reservations/{src}/{id}", reservationHandler.GetReservation)
		r.Post("/reservations/{src}/{id}", reservationHandler.PostReservation)
		r.Get("/reservations/{src}/{id}/process", reservationHandler.GetProcessReservation)
		r.Get("/reservations/{src}/{id}/delete", reservationHandler.GetDeleteReservation)
	})

	return mux
}
