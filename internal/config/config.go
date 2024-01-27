package config

import "github.com/go-webauthn/webauthn/webauthn"

// AppConfig holds the application config
type AppConfig struct {
	InProduction bool
	WebAuthn     *webauthn.WebAuthn
}
