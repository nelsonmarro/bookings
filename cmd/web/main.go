package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/background"
	"github.com/nelsonmarro/bookings/internal/driver"
	"github.com/nelsonmarro/bookings/internal/models"
	"github.com/nelsonmarro/bookings/internal/web"
)

const port = ":8080"

func main() {
	app := config.GetConfigInstance()
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	inProduction := flag.Bool("production", true, "Application is in production")
	dbName := flag.String("dbname", "bookingsdb", "Database name")
	dbUser := flag.String("dbuser", "nelson", "Database user")
	dbPassword := flag.String("dbpassword", "nelson9199", "Database password")
	dbPort := flag.String("dbport", "5432", "Database host")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbSSL := flag.String("dbssl", "disable", "Database SSL mode")

	flag.Parse()

	app.InProduction = *inProduction

	defer close(app.MailChan)

	fmt.Println("Starting mail listener...")
	background.ListenForMail(app)

	log.Println("Connecting to database...")

	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPassword, *dbSSL)
	db, err := driver.ConnectSQL(connectionString)
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
