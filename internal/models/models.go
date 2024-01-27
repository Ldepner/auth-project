package models

import (
	"github.com/go-webauthn/webauthn/webauthn"
	"time"
)

type TemplateData struct {
	Error     string
	Success   string
	StringMap map[string]string
	BoolMap   map[string]bool
}

type Session struct {
	Token                  string                `bson:"_id,omitempty"`
	UserID                 string                `bson:"user_id,omitempty"`
	Authenticated          bool                  `bson:"authenticated,omitempty"`
	LastActiveAt           time.Time             `bson:"last_active_at,omitempty"`
	ExpiresAt              time.Time             `bson:"expires_at,omitempty"`
	WebAuthNRegistration   *webauthn.SessionData `bson:"webauthn_registration,omitempty"`
	WebAuthNAuthentication *webauthn.SessionData `bson:"webauthn_authentication,omitempty"`
}
