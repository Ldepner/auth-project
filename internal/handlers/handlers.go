package handlers

import (
	"fmt"
	"github.com/Ldepner/auth-project/internal/authenticator"
	"github.com/Ldepner/auth-project/internal/config"
	"github.com/Ldepner/auth-project/internal/drivers"
	"github.com/Ldepner/auth-project/internal/helpers"
	"github.com/Ldepner/auth-project/internal/models"
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
	err := render.Template(w, "home.html", &models.TemplateData{})
	if err != nil {
		log.Println("error rendering home page")
		return
	}
}

// Login renders the login page
func (*Repository) Login(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, "login.html", &models.TemplateData{})
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

	loginForm := authenticator.NewLoginForm(r.Form)

	success, err := authenticator.Authenticate(loginForm)
	if err != nil {
		switch t := err.(type) {
		default:
			log.Println("error with login form", t)
			return
		case *helpers.ErrRecordNotFound:
			render.Template(w, "register.html", &models.TemplateData{Error: "email not found, please register"})
			return
		}
	}

	if !success {
		stringMap := make(map[string]string)
		stringMap["email"] = loginForm.Email
		err := render.Template(w, "login.html", &models.TemplateData{Error: "password incorrect", StringMap: stringMap})
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
	err := render.Template(w, "register.html", &models.TemplateData{})
	if err != nil {
		log.Println("error rendering registration page")
		return
	}
}

// PostRegister registers a new user
func (*Repository) PostRegister(w http.ResponseWriter, r *http.Request) {
	log.Println("attempting to register...")
	err := r.ParseForm()
	if err != nil {
		log.Println("error with registration form")
		return
	}

	regForm := authenticator.NewRegForm(r.PostForm)
	if regForm.Password != regForm.PasswordConfirmation {
		render.Template(w, "register.html", &models.TemplateData{})
		return
	}

	newRecord := models.UserRecord{
		Email:    regForm.Email,
		Password: regForm.Password,
	}

	err = Repo.DB.CreateUserRecord(&newRecord)
	if err != nil {
		switch t := err.(type) {
		default:
			log.Println(t)
			render.Template(w, "register.html", &models.TemplateData{})
			return
		case *helpers.ErrDuplicateEmail:
			render.Template(w, "login.html", &models.TemplateData{Error: "email already registered, please login"})
			return
		}
	}

	log.Println(fmt.Sprintf("user record creation successful: %s, %s", regForm.Email, regForm.Password))
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
