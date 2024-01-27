package models

import "github.com/go-webauthn/webauthn/webauthn"

type UserRecord struct {
	ID          string                `bson:"_id,omitempty"`
	Email       string                `bson:"email,omitempty"`
	Password    string                `bson:"password,omitempty"`
	Credentials []webauthn.Credential `bson:"credentials,omitempty"`
}

func (user *UserRecord) WebAuthnID() []byte {
	return []byte(user.ID)
}

func (user *UserRecord) WebAuthnName() string {
	return user.Email
}

func (user *UserRecord) WebAuthnDisplayName() string {
	return user.Email
}

func (user *UserRecord) WebAuthnIcon() string {
	return "https://pics.com/avatar.png"
}

func (user *UserRecord) WebAuthnCredentials() []webauthn.Credential {
	return user.Credentials
}
