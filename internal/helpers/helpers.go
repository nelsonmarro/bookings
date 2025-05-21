package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/nelsonmarro/bookings/config"
)

func ClientError(w http.ResponseWriter, status int) {
	app := config.GetConfigInstance()
	app.InfoLog.Println("Client error with status of:", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	app := config.GetConfigInstance()

	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
