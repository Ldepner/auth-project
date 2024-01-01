package dbrepo

import (
	"context"
	"github.com/Ldepner/auth-project/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func (r *mongoDBRepo) GetUserRecordByEmail(email string) (*models.UserRecord, error) {
	var user *models.UserRecord
	users := r.DB.Collection("users")

	filter := bson.D{{"email", email}}
	err := users.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *mongoDBRepo) CreateUserRecord(user *models.UserRecord) error {
	users := r.DB.Collection("users")

	// hash password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}
	newUserRecord := models.UserRecord{
		Email:    user.Email,
		Password: string(hashedBytes),
	}
	_, err = users.InsertOne(context.TODO(), newUserRecord)
	if err != nil {
		return err
	}

	return nil
}
