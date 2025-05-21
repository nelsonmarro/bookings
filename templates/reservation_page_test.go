package templates

import (
	"context"
	"io"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestReservationPage(t *testing.T) {
	vm := NewReservationPageVM("")
	r, w := io.Pipe()
	go func() {
		_ = ReservationPage(vm).Render(context.Background(), w)
		_ = w.Close()
	}()
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatalf("failed to create document: %v", err)
	}

	formElem := doc.Find(`[data-testid="reservationFormTest"]`)

	if formElem.Length() == 0 {
		t.Error("expected data-testid attribute to be rendered, but it wasn't")
	}

	if formElem.Find(`input[type="hidden"]`).Length() == 0 {
		t.Error("expected text input to be present")
	}

	if formElem.Find("#startdate").Length() == 0 {
		t.Error("expected start date input to be present")
	}
	if formElem.Find("#enddate").Length() == 0 {
		t.Error("expected end date input to be present")
	}

	if formElem.Find("button").Length() == 0 {
		t.Error("expected button to be present")
	}
}
