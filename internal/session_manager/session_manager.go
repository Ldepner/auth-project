package session_manager

import (
	"github.com/Ldepner/auth-project/internal/models"
	"github.com/Ldepner/auth-project/internal/repository"
	"net/http"
	"time"
)

var DBRepo repository.DBRepo

// NewSessionManager creates a new SessionManager
func NewSessionManager(dbRepo repository.DBRepo) {
	DBRepo = dbRepo
}

func AddSession(w http.ResponseWriter, userID string, expiresAt time.Time, authenticated bool) (string, error) {
	lastActiveAt := time.Now()
	token, err := DBRepo.CreateSession(userID, lastActiveAt, expiresAt, authenticated)
	if err != nil {
		return "", err
	}

	cookie := &http.Cookie{
		Name:     "id",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  expiresAt,
	}
	http.SetCookie(w, cookie)
	return token, nil
}

func AuthenticateSession(token string) error {
	err := DBRepo.UpdateSessionAuthenticated(token, true)
	return err
}

func InvalidateSession(w http.ResponseWriter, token string) error {
	err := DBRepo.DeleteSession(token)
	if err != nil {
		return err
	}

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

	return nil
}

func IsSessionTokenValid(token string) (bool, error) {
	session, err := DBRepo.GetSessionByID(token)

	if err != nil {
		return false, err
	}
	if !session.Authenticated {
		return false, nil
	}

	if time.Now().After(session.ExpiresAt) {
		_ = DBRepo.DeleteSession(token)
		return false, nil
	}

	return true, nil
}

func GetUserFromSession(token string) (*models.UserRecord, error) {
	session, err := DBRepo.GetSessionByID(token)
	if err != nil {
		return nil, err
	}

	user, err := DBRepo.GetUserRecordByID(session.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetSessionByID(token string) (*models.Session, error) {
	session, err := DBRepo.GetSessionByID(token)
	if err != nil {
		return nil, err
	}

	return session, nil
}
