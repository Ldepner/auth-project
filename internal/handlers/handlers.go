package handlers

import (
	"github.com/Ldepner/auth-project/internal/config"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

// NewRepo creates a new Repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home renders the home page
func (*Repository) Home(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("./templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Login renders the login page
func (*Repository) Login(w http.ResponseWriter, r *http.Request) {
	name, _ := filepath.Glob("./templates/login.html")
	tmpl, err := template.ParseFiles(name[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// PostLogin logs in a user
func (*Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("logging in")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// PostLogout signs a user out
func (*Repository) PostLogout(w http.ResponseWriter, r *http.Request) {
	log.Println("logging out")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Register renders the registration page
func (*Repository) Register(w http.ResponseWriter, r *http.Request) {
	name, _ := filepath.Glob("./templates/register.html")

	tmpl, err := template.ParseFiles(name[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// PostRegister registers a new user
func (*Repository) PostRegister(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
