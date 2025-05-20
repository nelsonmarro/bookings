package main

import (
	"encoding/gob"
	"log"
	"net/http"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/models"
	"github.com/nelsonmarro/bookings/internal/web"
)

const port = ":8080"

func main() {
	app := config.GetConfigInstance()
	gob.Register(models.Reservation{})

	err := http.ListenAndServe(port, web.Routes(app))
	if err != nil {
		log.Fatal(err)
	}
}
