package authenticator

import (
	"github.com/Ldepner/auth-project/internal/authenticator/auth_service"
	"github.com/Ldepner/auth-project/internal/repository"
	"net/url"
)

var DBRepo *repository.DBRepo

// NewAuthenticator creates a new Authenticator
func NewAuthenticator(dbRepo *repository.DBRepo) {
	DBRepo = dbRepo
}

type AuthService interface {
	Authenticate(data *LoginForm) (bool, error)
}

type LoginForm struct {
	Email    string
	Password string
}

type RegForm struct {
	Email                string
	Password             string
	PasswordConfirmation string
}

func NewLoginForm(data url.Values) *LoginForm {
	return &LoginForm{
		Email:    data.Get("email"),
		Password: data.Get("password"),
	}
}

func NewRegForm(data url.Values) *RegForm {
	return &RegForm{
		Email:                data.Get("email"),
		Password:             data.Get("password"),
		PasswordConfirmation: data.Get("confirmPassword"),
	}
}

func Authenticate(data *LoginForm) (string, bool, error) {
	// password strategy
	if len(data.Email) > 0 {
		authService := auth_service.PasswordAuthService{DBRepo: *DBRepo}
		userID, success, err := authService.Authenticate(data.Email, data.Password)

		return userID, success, err
	}

	return "", false, nil
}
