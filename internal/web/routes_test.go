package web

import (
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/nelsonmarro/bookings/config"
)

func TestRoutes(t *testing.T) {
	app := config.GetConfigInstance()
	mux := Routes(app)

	switch v := mux.(type) {
	case *chi.Mux:
	default:
		t.Errorf("Expected mux to be of type *chi.Mux, got %T", v)
	}
}
