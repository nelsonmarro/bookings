package dbrepo

import (
	"database/sql"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/repository"
)

type postgresDbRepo struct {
	DB  *sql.DB
	App *config.AppConfig
}

type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewTestingRepo(a *config.AppConfig) repository.DataBaseRepo {
	return &testDBRepo{
		App: a,
	}
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DataBaseRepo {
	return &postgresDbRepo{
		DB:  conn,
		App: a,
	}
}
