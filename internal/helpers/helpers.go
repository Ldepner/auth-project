package helpers

import (
	"github.com/Ldepner/auth-project/internal/config"
	"net/http"
)

var app *config.AppConfig

// NewHelpers sets up appConfig for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

func IsAuthenticated(r *http.Request) bool {
	//exists := app.Session.Exists(r.Context(), "user_id")
	exists, err := r.Cookie("id")
	if err != nil {
		return false
	}

	if len(exists.Value) > 0 {
		return true
	}

	return false
}
