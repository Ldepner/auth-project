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
	"github.com/Ldepner/auth-project/internal/session_manager"
	"log"
	"net/http"
	"time"
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
	// get user from session
	sessionToken, err := r.Cookie("id")
	if err != nil {
		response := map[string]interface{}{
			"message": "error finding session. please login.",
		}
		JSONResponse(w, 400, response)
		return
	}
	user, err := session_manager.GetUserFromSession(sessionToken.Value)

	// does user have webauthN credentials already?
	biometricRegistered := len(user.Credentials) > 0
	data := make(map[string]bool)
	data["biometricRegistered"] = biometricRegistered

	err = render.Template(w, "home.html", &models.TemplateData{BoolMap: data})
	if err != nil {
		log.Println("error rendering home page")
		return
	}
}

// Login renders the login page
func (*Repository) Login(w http.ResponseWriter, r *http.Request) {
	if helpers.IsAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

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

	userID, success, err := authenticator.Authenticate(loginForm)
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
	_, err = session_manager.AddSession(w, userID, time.Now().Local().Add(30*time.Minute), true)
	if err != nil {
		render.Template(w, "login.html", &models.TemplateData{Error: "error occured, please try again"})
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// PostLogout signs a user out
func (*Repository) PostLogout(w http.ResponseWriter, r *http.Request) {
	log.Println("logging out")
	token, err := r.Cookie("id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	_ = session_manager.InvalidateSession(w, token.Value)
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

	userID, err := Repo.DB.CreateUserRecord(&newRecord)
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
	_, err = session_manager.AddSession(w, userID, time.Now().Local().Add(30*time.Minute), true)
	if err != nil {
		render.Template(w, "register.html", &models.TemplateData{Error: "an error occured, please try again"})
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// BeginRegistration begins registering a user with webauthn
func (*Repository) BeginRegistration(w http.ResponseWriter, r *http.Request) {
	// Get user from session
	sessionToken, err := r.Cookie("id")
	if err != nil {
		response := map[string]interface{}{
			"message": "error finding session. please login.",
		}
		JSONResponse(w, 400, response)
		return
	}
	user, err := session_manager.GetUserFromSession(sessionToken.Value)
	if err != nil {
		response := map[string]interface{}{
			"message": "error finding user",
		}
		JSONResponse(w, 400, response)
		return
	}

	// generate PublicKeyCredentialCreationOptions, session data
	options, sessionData, err := Repo.App.WebAuthn.BeginRegistration(user)

	if err != nil {
		log.Println(err)
		JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{"message": err.Error()})
		return
	}

	// store session data as marshaled JSON
	err = Repo.DB.UpdateWebAuthNSession(sessionToken.Value, "webauthn_registration", sessionData)
	if err != nil {
		log.Println(err)
		JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{"message": err.Error()})
		return
	}

	JSONResponse(w, http.StatusOK, options)
}

// FinishRegistration finishes registering a user with webauthn
func (*Repository) FinishRegistration(w http.ResponseWriter, r *http.Request) {
	log.Println("attempting to finish webauthn registration...")
	sessionToken, err := r.Cookie("id")
	if err != nil {
		response := map[string]interface{}{
			"message": "error finding session. please login.",
		}
		JSONResponse(w, 400, response)
		return
	}

	// Get user from session
	user, err := session_manager.GetUserFromSession(sessionToken.Value)
	if err != nil {
		response := map[string]interface{}{
			"message": "error finding user",
		}
		JSONResponse(w, 400, response)
		return
	}

	// Get session including webauthn session data
	session, err := session_manager.GetSessionByID(sessionToken.Value)
	if err != nil {
		response := map[string]interface{}{
			"message": "error finding session data for registration",
		}
		JSONResponse(w, 400, response)
		return
	}

	sessionData := *session.WebAuthNRegistration

	// Finish registration
	credential, err := Repo.App.WebAuthn.FinishRegistration(user, sessionData, r)
	if err != nil {
		log.Println(err)
		JSONResponse(w, http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
		return
	}

	// Add Credential to user
	err = Repo.DB.UpdateUserRecord(user.ID, "credentials", append(user.Credentials, *credential))
	if err != nil {
		log.Println(err)
		JSONResponse(w, http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
		return
	}

	log.Println("registration success")

	response := map[string]interface{}{
		"message": "Registration success",
	}
	JSONResponse(w, 200, response)
	return
}

// BeginLogin begins logging a user in with webauthn
func (*Repository) BeginLogin(w http.ResponseWriter, r *http.Request) {
	// get email
	log.Println("beginning login...")
	email := r.URL.Query().Get("email")
	if email == "" {
		log.Println("error with login page")
		body := map[string]interface{}{
			"message": "error with login form",
		}
		JSONResponse(w, http.StatusBadRequest, body)
		return
	}

	// get user
	log.Println(fmt.Sprintf("finding user by email %s", email))
	user, err := Repo.DB.GetUserRecordByEmail(email)
	if err != nil {
		log.Println("user not found, please register or try again")
		body := map[string]interface{}{
			"message": "error with login form",
		}
		JSONResponse(w, http.StatusBadRequest, body)
		return
	}

	// generate PublicKeyCredentialRequestOptions, session data
	options, sessionData, err := Repo.App.WebAuthn.BeginLogin(user)
	if err != nil {
		log.Println("error generating login session data")
		body := map[string]interface{}{
			"message": err.Error(),
		}
		JSONResponse(w, http.StatusBadRequest, body)
		return
	}
	// add session with expires at in past, so will not allow authentication
	token, _ := session_manager.AddSession(w, user.ID, time.Now().Add(30*time.Minute), false)

	// store session data as marshaled JSON
	err = Repo.DB.UpdateWebAuthNSession(token, "webauthn_authentication", sessionData)
	if err != nil {
		log.Println("error updating session")
		body := map[string]interface{}{
			"message": err.Error(),
		}
		JSONResponse(w, http.StatusBadRequest, body)
		return
	}

	JSONResponse(w, http.StatusOK, options)
}

// FinishLogin finishes logging in a user with webauthn
func (*Repository) FinishLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("attempting to finish webauthn login...")
	sessionToken, err := r.Cookie("id")
	if err != nil {
		response := map[string]interface{}{
			"message": "error finding session. please login.",
		}
		JSONResponse(w, 400, response)
		return
	}

	// Get user from session
	user, err := session_manager.GetUserFromSession(sessionToken.Value)
	if err != nil {
		response := map[string]interface{}{
			"message": "error finding user",
		}
		JSONResponse(w, 400, response)
		return
	}

	// Get session including webauthn session data
	session, err := session_manager.GetSessionByID(sessionToken.Value)
	if err != nil {
		response := map[string]interface{}{
			"message": "error finding session data for login",
		}
		JSONResponse(w, 400, response)
		return
	}

	// in an actual implementation we should perform additional
	// checks on the returned 'credential'
	sessionData := *session.WebAuthNAuthentication

	_, err = Repo.App.WebAuthn.FinishLogin(user, sessionData, r)
	if err != nil {
		log.Println(err)
		response := map[string]interface{}{
			"message": err.Error(),
		}
		JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Add new valid session
	log.Println(fmt.Sprintf("authentication success"))
	err = session_manager.AuthenticateSession(sessionToken.Value)
	if err != nil {
		log.Println(err)
		response := map[string]interface{}{
			"message": err.Error(),
		}
		JSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// handle successful login
	response := map[string]interface{}{
		"message": "success",
	}
	JSONResponse(w, http.StatusOK, response)
}
