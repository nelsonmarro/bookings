package handlers

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/email"
	"github.com/nelsonmarro/bookings/internal/handlers/rooms"
	"github.com/nelsonmarro/bookings/internal/middlewares"
	"github.com/nelsonmarro/bookings/internal/models"
	"github.com/nelsonmarro/bookings/internal/repository/dbrepo"
)

func TestMain(m *testing.M) {
	app := config.GetConfigInstance()

	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	mailChan := make(chan email.MailData)
	app.MailChan = mailChan
	defer close(app.MailChan)

	lisentForMail(app)

	os.Exit(m.Run())
}

func lisentForMail(app *config.AppConfig) {
	go func() {
		for {
			_ = <-app.MailChan
		}
	}()
}

func getRoutes() http.Handler {
	app := config.GetConfigInstance()
	dbrepo := dbrepo.NewTestingRepo(app)

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

	homepageHandler := NewHomepageHandler(app)
	mux.Get("/", homepageHandler.Get)

	aboutpageHandler := NewAboutpageHandler(app)
	mux.Get("/about", aboutpageHandler.Get)

	contactpageHandler := NewContactpageHandler(app)
	mux.Get("/contact", contactpageHandler.Get)

	reservationPageHandler := NewReservationpageHandler(app, dbrepo)
	mux.Get("/reservation", reservationPageHandler.Get)
	mux.Post("/reservation", reservationPageHandler.Post)
	mux.Post("/reservation-json", reservationPageHandler.PostJson)

	chooseRoomHandler := NewChooseRoomHandler(app, dbrepo)
	mux.Get("/choose-room/{id}", chooseRoomHandler.Get)

	bookRoomHandler := NewBookRoomHandler(app, dbrepo)
	mux.Get("/book-room", bookRoomHandler.Get)

	singleRoomHandler := rooms.NewSingleRoomHandler(app)
	mux.Get("/rooms/single", singleRoomHandler.Get)

	confirmReservationHandler := NewConfirmReservationHandler(app, dbrepo)
	mux.Get("/reservation/confirmation", confirmReservationHandler.Get)
	mux.Post("/reservation/confirmation", confirmReservationHandler.Post)

	reservationSummaryHandler := NewReservationSummaryHandler(app)
	mux.Get("/reservation/summary", reservationSummaryHandler.Get)

	doubleRoomHandler := rooms.NewDoubleRoomHandler(app)
	mux.Get("/rooms/double", doubleRoomHandler.Get)

	return mux
}
