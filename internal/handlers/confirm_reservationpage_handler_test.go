package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"github.com/nelsonmarro/bookings/tests"
)

func TestPostConfirmReservationTempl(t *testing.T) {
	confirmationTests := []struct {
		name               string
		url                string
		expectedStatusCode int
		params             []tests.PostData
		assertDoc          func(doc *goquery.Document)
	}{
		{
			name:               "post-confirmation-whith-empty-fields",
			url:                "/reservation/confirmation",
			expectedStatusCode: http.StatusOK,
			params:             []tests.PostData{},
			assertDoc: func(doc *goquery.Document) {
				expectedText := "This field it's requiredThis field it's requiredThis field it's required"

				formElement := doc.Find(`[data-testid="confirm-reservation-form"]`)
				if formElement.Length() == 0 {
					t.Error("expected form element to be present")
				}
				pErrors := formElement.Find("p")
				if pErrors.Length() == 0 {
					t.Error("expected p element to be present")
				}
				if pErrors.Text() != expectedText {
					t.Errorf("expected p text to be '%s', got '%s'", expectedText, pErrors.Text())
				}
			},
		},
		{
			name:               "post-confirmation-whith-all-fields",
			url:                "/reservation/confirmation",
			expectedStatusCode: http.StatusOK,
			params: []tests.PostData{
				{Key: "first_name", Value: "Nelson"},
				{Key: "last_name", Value: "Marro"},
				{Key: "email", Value: "nelson@gmail.com"},
				{Key: "phone_number", Value: "123456789"},
			},
			assertDoc: func(doc *goquery.Document) {
				formElement := doc.Find(`[data-testid="confirm-reservation-form"]`)
				if formElement.Length() != 0 {
					t.Error("expected form element to be absent")
				}
			},
		},
	}

	mux := getRoutes()
	ts := httptest.NewTLSServer(mux)
	defer ts.Close()

	for _, e := range confirmationTests {
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

		doc, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			t.Fatalf("failed to read template: %v", err)
		}

		e.assertDoc(doc)
	}
}
