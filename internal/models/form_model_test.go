package models

import (
	"net/http"
	"net/url"
	"testing"
)

func newTestForm(values url.Values) *Form {
	return &Form{
		Values: values,
		Errors: make(Errors),
	}
}

func TestForm_Required(t *testing.T) {
	form := newTestForm(url.Values{
		"first_name": {"John"},
		"last_name":  {""},
		"email":      {""},
	})

	form.Required("first_name", "last_name", "email")

	if _, ok := form.Errors["email"]; !ok {
		t.Error("Expected error for required email field")
	}
	if _, ok := form.Errors["last_name"]; !ok {
		t.Error("Expected error for required last_name field")
	}
	if _, ok := form.Errors["first_name"]; ok {
		t.Error("Did not expect error for filled first_name field")
	}
}

func TestForm_Has(t *testing.T) {
	form := newTestForm(url.Values{})
	req := &http.Request{Form: url.Values{"first_name": {"John"}}}

	if !form.Has("first_name", req) {
		t.Error("Expected Has to return true for present field")
	}
	if form.Has("baz", req) {
		t.Error("Expected Has to return false for missing field")
	}
}

func TestForm_MinLength(t *testing.T) {
	form := newTestForm(url.Values{})
	req := &http.Request{Form: url.Values{"password": {"abc"}}}

	ok := form.MinLength("password", 5, req)
	if ok {
		t.Error("Expected MinLength to return false for short value")
	}
	if _, ok := form.Errors["password"]; !ok {
		t.Error("Expected error for short password")
	}

	req = &http.Request{Form: url.Values{"password": {"abcdef"}}}
	form.Errors = Errors{}
	ok = form.MinLength("password", 5, req)
	if !ok {
		t.Error("Expected MinLength to return true for long enough value")
	}
}

func TestForm_IsEmail(t *testing.T) {
	form := newTestForm(url.Values{"email": {"invalid-email"}})
	if form.IsEmail("email") {
		t.Error("Expected IsEmail to return false for invalid email")
	}
	if _, ok := form.Errors["email"]; !ok {
		t.Error("Expected error for invalid email")
	}

	form = newTestForm(url.Values{"email": {"test@example.com"}})
	form.Errors = Errors{}
	if !form.IsEmail("email") {
		t.Error("Expected IsEmail to return true for valid email")
	}
}

func TestForm_Valid(t *testing.T) {
	form := newTestForm(url.Values{})
	if !form.Valid() {
		t.Error("Expected Valid to return true when no errors")
	}
	form.Errors.Add("foo", "bar")
	if form.Valid() {
		t.Error("Expected Valid to return false when errors exist")
	}
}
