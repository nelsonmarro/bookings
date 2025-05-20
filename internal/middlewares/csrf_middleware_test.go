package middlewares

import (
	"net/http"
	"testing"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/tests"
)

func TestCSRFMiddleware(t *testing.T) {
	var mh tests.MyHandler
	app := config.GetConfigInstance()
	h := CSRFMiddleware(&mh, app)

	switch v := h.(type) {
	case http.Handler:
	default:
		t.Errorf("Expected http.Handler, got %T", v)
	}
}
