package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/models"
	"github.com/nelsonmarro/bookings/internal/repository/dbrepo"
)

var testCases = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"single-room", "/rooms/single", "GET", http.StatusOK},
	{"double-room", "/rooms/double", "GET", http.StatusOK},
}

func TestHandlers(t *testing.T) {
	mux := getRoutes()
	ts := httptest.NewTLSServer(mux)
	defer ts.Close()

	for _, e := range testCases {
		if e.method == "GET" {
			response, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Errorf("Error making GET request: %v", err)
			}

			if response.StatusCode != e.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", e.expectedStatusCode, response.StatusCode)
			}
		}
	}
}

func TestRepository_Post_ConfirmReservation(t *testing.T) {
	reqBody := "start_date=2025-10-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2025-10-05")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Doe")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=jhon@email.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone_number=0987654321")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/reservation/confirmation", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	app := config.GetConfigInstance()
	dbrepo := dbrepo.NewTestingRepo(app)

	confirmReservationHandler := NewConfirmReservationHandler(app, dbrepo)

	handler := http.HandlerFunc(confirmReservationHandler.Post)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Expected status code %d, got %d", http.StatusSeeOther, rr.Code)
	}

	// test for missing post body
	req, _ = http.NewRequest("POST", "/reservation/confirmation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	confirmReservationHandler = NewConfirmReservationHandler(app, dbrepo)

	handler = http.HandlerFunc(confirmReservationHandler.Post)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected status code %d, got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// test for invalid start date
	reqBody = "start_date=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2025-10-05")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Doe")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=jhon@email.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone_number=0987654321")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/reservation/confirmation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	confirmReservationHandler = NewConfirmReservationHandler(app, dbrepo)

	handler = http.HandlerFunc(confirmReservationHandler.Post)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected status code %d, got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// test for faliure insert reservation
	reqBody = "start_date=2025-10-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2025-10-05")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Doe")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=jhon@email.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone_number=0987654321")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")

	req, _ = http.NewRequest("POST", "/reservation/confirmation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	confirmReservationHandler = NewConfirmReservationHandler(app, dbrepo)

	handler = http.HandlerFunc(confirmReservationHandler.Post)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Handler failed when trying inserting reservation. Expected status code %d, got %d", http.StatusTemporaryRedirect, rr.Code)
	}
}

func TestRepository_Get_ConfirmReservation(t *testing.T) {
	app := config.GetConfigInstance()
	dbrepo := dbrepo.NewTestingRepo(app)

	res := models.Reservation{
		RoomID: 1,
		Room:   models.Room{ID: 1, RoomName: "Single Room"},
	}

	req, _ := http.NewRequest("GET", "/reservation/confirmation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	app.Session.Put(ctx, "reservation", res)

	confirmReservationHandler := NewConfirmReservationHandler(app, dbrepo)
	handler := http.HandlerFunc(confirmReservationHandler.Get)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	req, _ = http.NewRequest("GET", "/reservation/confirmation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected status code %d, got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// test with no existing room
	req, _ = http.NewRequest("GET", "/reservation/confirmation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	res.RoomID = 100
	app.Session.Put(ctx, "reservation", res)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}
}

func getCtx(req *http.Request) context.Context {
	app := config.GetConfigInstance()
	ctx, err := app.Session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Printf("Error loading session: %v", err)
	}

	return ctx
}
