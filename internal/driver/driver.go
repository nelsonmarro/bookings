package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// DB holds a database connection.
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const (
	maxOpenDBConn = 10
	maxIdleDBConn = 5
	maxDBLifeTime = 5 * time.Minute
)

// Creates database pool for postgres
func ConnectSQL(dsn string) (*DB, error) {
	db, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(maxOpenDBConn)
	db.SetMaxIdleConns(maxIdleDBConn)
	db.SetConnMaxLifetime(maxDBLifeTime)

	dbConn.SQL = db

	err = testDB(db)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

// Tries to Ping the DB
func testDB(db *sql.DB) error {
	err := db.Ping()
	return err
}

// Creates a new Database for the application
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
