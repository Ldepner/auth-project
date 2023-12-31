package authenticator

import (
	"github.com/Ldepner/auth-project/internal/authenticator/auth_service"
	"github.com/Ldepner/auth-project/internal/repository"
	"log"
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

func NewLoginForm(data url.Values) *LoginForm {
	return &LoginForm{
		Email:    data.Get("email"),
		Password: data.Get("password"),
	}
}

func Authenticate(data *LoginForm) (bool, error) {
	// password strategy
	log.Println(data.Email)
	if len(data.Email) > 0 {
		log.Println("test")
		authService := auth_service.PasswordAuthService{}
		success, err := authService.Authenticate(data.Email, data.Password)

		return success, err
	}

	return false, nil
}
