package handlers

import (
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/helpers"
	"github.com/nelsonmarro/bookings/internal/models"
	"github.com/nelsonmarro/bookings/internal/repository"
	"github.com/nelsonmarro/bookings/templates/user"
)

type UserHandler struct {
	app *config.AppConfig
	DB  repository.DataBaseRepo
}

func NewUserHandler(app *config.AppConfig, dbrepo repository.DataBaseRepo) *UserHandler {
	return &UserHandler{
		app: app,
		DB:  dbrepo,
	}
}

func (h *UserHandler) GetLogin(w http.ResponseWriter, r *http.Request) {
	vm := user.NewLoginPageVM()
	vm.CSRFToken = nosurf.Token(r)
	loginPage := user.LoginPage(vm)
	err := loginPage.Render(r.Context(), w)
	if err != nil {
		helpers.ServerError(w, err)
	}
}

func (h *UserHandler) PostLogin(w http.ResponseWriter, r *http.Request) {
	// renew session token
	_ = h.app.Session.RenewToken(r.Context())

	if err := r.ParseForm(); err != nil {
		h.app.ErrorLog.Printf("error parsing form: %v", err)
	}

	form := models.NewForm(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	vm := user.NewLoginPageVM()
	if !form.Valid() {
		vm.Form = form
		vm.CSRFToken = nosurf.Token(r)
		loginPage := user.LoginPage(vm)
		err := loginPage.Render(r.Context(), w)
		if err != nil {
			h.app.ErrorLog.Printf("error rendering login page: %v", err)
			return
		}
	}

	email := form.Get("email")
	password := form.Get("password")

	id, _, err := h.DB.Authenticate(email, password)
	if err != nil {
		h.app.ErrorLog.Printf("error authenticating user: %v", err)
		h.app.Session.Put(r.Context(), "error", "invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	h.app.Session.Remove(r.Context(), "error")
	h.app.Session.Put(r.Context(), "user_id", id)
	h.app.Session.Put(r.Context(), "info", "logged in successfully")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout handler to log out the user
func (h *UserHandler) GetLogout(w http.ResponseWriter, r *http.Request) {
	h.app.Session.Destroy(r.Context())
	h.app.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
