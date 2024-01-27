package helpers

import (
	"github.com/Ldepner/auth-project/internal/config"
	"github.com/Ldepner/auth-project/internal/session_manager"
	"net/http"
)

var app *config.AppConfig

// NewHelpers sets up appConfig for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

func IsAuthenticated(r *http.Request) bool {
	exists, err := r.Cookie("id")
	if err != nil {
		return false
	}

	// Check if session is valid
	isValid, err := session_manager.IsSessionTokenValid(exists.Value)
	if err != nil {
		return false
	}

	return isValid
}
