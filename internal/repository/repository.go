package repository

import (
	"github.com/Ldepner/auth-project/internal/models"
	"github.com/go-webauthn/webauthn/webauthn"
	"time"
)

type DBRepo interface {
	GetUserRecordByEmail(email string) (*models.UserRecord, error)
	CreateUserRecord(user *models.UserRecord) (string, error)
	CreateSession(userID string, lastActiveAt, expiresAt time.Time, authenticated bool) (string, error)
	DeleteSession(token string) error
	GetSessionByID(token string) (*models.Session, error)
	GetUserRecordByID(id string) (*models.UserRecord, error)
	UpdateSessionAuthenticated(token string, authenticated bool) error
	UpdateWebAuthNSession(token string, sessionType string, webauthnSessionData *webauthn.SessionData) error
	UpdateUserRecord(userID string, field string, updatedValue any) error
}
