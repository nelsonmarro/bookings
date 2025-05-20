package models

import (
	"net/http"

	"github.com/nelsonmarro/bookings/config"
)

type Errors map[string][]string

func (e Errors) Add(field, message string) {
	if e == nil {
		e = make(Errors)
	}
	e[field] = append(e[field], message)
}

func (e Errors) Get(field string) string {
	es := e[field]
	return es[0]
}

func (e Errors) HasField(field string) bool {
	_, exists := e[field]
	return exists
}

type MessageType string

const (
	MessageTypeError   MessageType = "error"
	MessageTypeInfo    MessageType = "info"
	MessageTypeWarning MessageType = "warning"
)

type BaseViewModel struct {
	Form        *Form
	FormErrors  Errors
	CSRFToken   string
	MessageType MessageType
	Message     string
}

func GetSessionMessage(app *config.AppConfig, r *http.Request) (MessageType, string) {
	errorMessage := app.Session.PopString(r.Context(), "error")
	if errorMessage != "" {
		return MessageTypeError, errorMessage
	}
	infoMessage := app.Session.PopString(r.Context(), "info")
	if infoMessage != "" {
		return MessageTypeInfo, infoMessage
	}
	warningMessage := app.Session.PopString(r.Context(), "warning")
	if warningMessage != "" {
		return MessageTypeWarning, warningMessage
	}

	return "", ""
}
