package main

import (
	"net/http"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/web"
)

const port = ":8080"

func main() {
	app := config.GetConfigInstance()

	_ = http.ListenAndServe(port, web.Routes(app))
}
