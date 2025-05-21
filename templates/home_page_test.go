package templates

import (
	"context"
	"io"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestHomePage(t *testing.T) {
	vm := NewHomePageVM()
	r, w := io.Pipe()
	go func() {
		_ = HomePage(vm).Render(context.Background(), w)
		_ = w.Close()
	}()
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatalf("failed to create document: %v", err)
	}

	if doc.Find(`[data-testid="headerTempl"]`).Length() == 0 {
		t.Errorf("expected data-testid attribute to be present")
	}

	// Expect the page name to be set correctly
	expectedPageName := "Home Page"
	if doc.Find("title").Text() != expectedPageName {
		t.Errorf("expected page name to be '%s', got '%s'", expectedPageName, doc.Find("title").Text())
	}

	// Expect the component to include a testid.
	if doc.Find(`[data-testid="navTempl"]`).Length() == 0 {
		t.Error("expected data-testid attribute to be rendered, but it wasn't")
	}

	// Expect the component to include a testid.
	if doc.Find(`[data-testid="homeTempl"]`).Length() == 0 {
		t.Error("expected data-testid attribute to be rendered, but it wasn't")
	}

	// Expect h1 element to be present.
	if doc.Find("h1").Length() == 0 {
		t.Error("expected h1 element to be present")
	}
	if doc.Find("h1").Text() != "Welcome to Our Booking Site" {
		t.Errorf("expected h1 text to be 'Welcome to Our Booking Site', got '%s'", doc.Find("h1").Text())
	}

	// Expect the component to include a testid.
	if doc.Find(`[data-testid="footerTempl"]`).Length() == 0 {
		t.Error("expected data-testid attribute to be rendered, but it wasn't")
	}

	// Expect the footer to have a div whith 3 div child elements.
	if doc.Find("footer").Length() == 0 {
		t.Error("expected footer element to be present")
	}

	if doc.Find("footer > div > div").Length() != 3 {
		t.Errorf("expected footer to have 3 child div elements, got %d", doc.Find("footer > div > div").Length())
	}
}
