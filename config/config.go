package config

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/alexedwards/scs/v2"
)

var (
	once     sync.Once
	instance *AppConfig
)

type AppConfig struct {
	UseCache     bool
	InfoLog      *log.Logger
	ErrorLog     *log.Logger
	InProduction bool
	Session      *scs.SessionManager
}

func GetConfigInstance() *AppConfig {
	once.Do(func() {
		instance = &AppConfig{
			Session:      scs.New(),
			UseCache:     false,
			InfoLog:      log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
			ErrorLog:     log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
			InProduction: false,
		}
		instance.Session.Lifetime = 24 * time.Hour // 24 horas
		instance.Session.Cookie.Persist = true
		instance.Session.Cookie.SameSite = http.SameSiteLaxMode
		instance.Session.Cookie.Secure = instance.InProduction
	})
	return instance
}
