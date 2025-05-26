package web

import (
	"log"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/driver"
)

func TestRoutes(t *testing.T) {
	app := config.GetConfigInstance()
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookingsdb user=nelsonmarro password=nelson9199 sslmode=disable")
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	log.Println("Connected to database...")
	defer db.SQL.Close()

	mux := Routes(app, db)

	switch v := mux.(type) {
	case *chi.Mux:
	default:
		t.Errorf("Expected mux to be of type *chi.Mux, got %T", v)
	}
}
