package test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/nelsonmarro/bookings/tests"
)

var testCases = []struct {
	name               string
	url                string
	method             string
	params             []tests.PostData
	expectedStatusCode int
}{
	{"home", "/", "GET", []tests.PostData{}, http.StatusOK},
	{"about", "/about", "GET", []tests.PostData{}, http.StatusOK},
	{"contact", "/contact", "GET", []tests.PostData{}, http.StatusOK},
	{"single-room", "/rooms/single", "GET", []tests.PostData{}, http.StatusOK},
	{"double-room", "/rooms/double", "GET", []tests.PostData{}, http.StatusOK},
	{"resrvation-search", "/reservation", "GET", []tests.PostData{}, http.StatusOK},
	{"resrvation-confirmation", "/reservation/confirmation", "GET", []tests.PostData{}, http.StatusOK},
	{"resrvation-summary", "/reservation/summary", "GET", []tests.PostData{}, http.StatusOK},
	{"post-reservation-search", "/reservation", "POST", []tests.PostData{
		{Key: "startdate", Value: "2023-10-01"},
		{Key: "enddate", Value: "2023-10-09"},
	}, http.StatusOK},
	{"post-reservation-search-json", "/reservation-json", "POST", []tests.PostData{
		{Key: "startdate", Value: "2023-10-01"},
		{Key: "enddate", Value: "2023-10-09"},
	}, http.StatusOK},
	{"confirm-reservation", "/reservation-json", "POST", []tests.PostData{
		{Key: "first_name", Value: "Nelson"},
		{Key: "last_name", Value: "Marro"},
		{Key: "email", Value: "nelson@gmail.com"},
		{Key: "phone_number", Value: "123456789"},
	}, http.StatusOK},
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
		} else {
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.Key, x.Value)
			}

			response, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Errorf("Error making POST request: %v", err)
			}
			if response.StatusCode != e.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", e.expectedStatusCode, response.StatusCode)
			}
		}
	}
}
