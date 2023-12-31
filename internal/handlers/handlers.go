package handlers

import (
	"fmt"
	"github.com/Ldepner/auth-project/internal/authenticator"
	"github.com/Ldepner/auth-project/internal/config"
	"github.com/Ldepner/auth-project/internal/drivers"
	"github.com/Ldepner/auth-project/internal/render"
	"github.com/Ldepner/auth-project/internal/repository"
	"github.com/Ldepner/auth-project/internal/repository/dbrepo"
	"log"
	"net/http"
)

type Repository struct {
	App *config.AppConfig
	DB  repository.DBRepo
}

var Repo *Repository

// NewRepo creates a new Repository
func NewRepo(a *config.AppConfig, db *drivers.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewMongoRepo(db.NoSQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home renders the home page
func (*Repository) Home(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, "home.html")
	if err != nil {
		log.Println("error rendering home page")
		return
	}
}

// Login renders the login page
func (*Repository) Login(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, "login.html")
	if err != nil {
		log.Println("error rendering login page")
		return
	}
}

// PostLogin logs in a user
func (*Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("attempting to log in...")
	err := r.ParseForm()
	if err != nil {
		log.Println("error with login page")
		return
	}

	log.Println(r.PostForm)
	loginForm := authenticator.NewLoginForm(r.Form)

	success, err := authenticator.Authenticate(loginForm)
	if err != nil {
		log.Println("error with login form")
		return
	}
	log.Println(success == false)

	if !success {
		// TODO: Add flash for unsuccessful login
		err := render.Template(w, "login.html")
		if err != nil {
			log.Println("error rendering login page")
			return
		}
		return
	}

	log.Println(fmt.Sprintf("authentication success: %t", success))
	cookie := &http.Cookie{
		Name:     "id",
		Value:    "1",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// PostLogout signs a user out
func (*Repository) PostLogout(w http.ResponseWriter, r *http.Request) {
	log.Println("logging out")
	cookie := &http.Cookie{
		Name:     "id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Register renders the registration page
func (*Repository) Register(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, "register.html")
	if err != nil {
		log.Println("error rendering registration page")
		return
	}
}

// PostRegister registers a new user
func (*Repository) PostRegister(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
