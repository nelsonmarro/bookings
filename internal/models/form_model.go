package models

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-playground/validator/v10"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

type Form struct {
	url.Values
	Errors Errors
}

func NewForm(values url.Values) *Form {
	return &Form{
		Values: values,
		Errors: make(Errors),
	}
}

func (f *Form) Required(fileds ...string) {
	for _, field := range fileds {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field it's required")
		}
	}
}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.FormValue(field)
	return x != ""
}

func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	x := r.FormValue(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

func (f *Form) IsEmail(field string) bool {
	validate = validator.New()
	err := validate.Var(f.Get("email"), "email")
	if err != nil {
		f.Errors.Add(field, "the field its not an valid email")
		return false
	}
	return true
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
