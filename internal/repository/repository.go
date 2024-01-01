package repository

import "github.com/Ldepner/auth-project/internal/models"

type DBRepo interface {
	GetUserRecordByEmail(email string) (*models.UserRecord, error)
	CreateUserRecord(user *models.UserRecord) error
}
