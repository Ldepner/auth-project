package auth_service

import (
	"github.com/Ldepner/auth-project/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type PasswordAuthService struct {
	DBRepo repository.DBRepo
}

func (p *PasswordAuthService) Authenticate(email, password string) (bool, error) {
	user, err := p.DBRepo.GetUserRecordByEmail(email)
	if err != nil {
		log.Println(err)
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
