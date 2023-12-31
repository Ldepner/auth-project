package auth_service

import (
	"errors"
	"github.com/Ldepner/auth-project/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type PasswordAuthService struct {
	DBRepo repository.DBRepo
}

func (p *PasswordAuthService) Authenticate(email, password string) (bool, error) {
	user, err := p.DBRepo.GetUserRecordByEmail(email)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, errors.New("incorrect password")
	} else if err != nil {
		return false, err
	}

	return true, nil
}
