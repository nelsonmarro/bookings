package main

import (
	"encoding/gob"
	"log"
	"net/http"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/driver"
	"github.com/nelsonmarro/bookings/internal/models"
	"github.com/nelsonmarro/bookings/internal/web"
)

const port = ":8080"

func main() {
	app := config.GetConfigInstance()
	gob.Register(models.Reservation{})

	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookingsdb user=nelson password=nelson9199 sslmode=disable")
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	log.Println("Connected to database...")
	defer db.SQL.Close()

	err = http.ListenAndServe(port, web.Routes(app, db))
	if err != nil {
		log.Fatal(err)
	}
}
